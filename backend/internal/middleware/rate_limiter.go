package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimitStrategy defines the rate limiting strategy
type RateLimitStrategy string

const (
	// StrategyIP - Rate limit by IP address
	StrategyIP RateLimitStrategy = "ip"
	// StrategyUser - Rate limit by user ID (authenticated users)
	StrategyUser RateLimitStrategy = "user"
	// StrategyIPAndUser - Rate limit by both IP and user ID
	StrategyIPAndUser RateLimitStrategy = "ip_and_user"
	// StrategyEndpoint - Rate limit by endpoint
	StrategyEndpoint RateLimitStrategy = "endpoint"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	// Enabled - Enable/disable rate limiting
	Enabled bool
	// Strategy - Rate limiting strategy (ip, user, ip_and_user, endpoint)
	Strategy RateLimitStrategy
	// RequestsPerSecond - Number of requests allowed per second
	RequestsPerSecond int
	// Burst - Maximum burst size
	Burst int
	// Redis - Redis client for distributed rate limiting (nil = in-memory)
	Redis *redis.Client
	// Logger - Logger instance
	Logger *zap.Logger
	// SkipSuccessful - Skip rate limiting for successful requests (for testing)
	SkipSuccessful bool
	// CustomKeyFunc - Custom function to generate rate limit key
	CustomKeyFunc func(c echo.Context) string
	// OnLimitReached - Callback when rate limit is reached
	OnLimitReached func(c echo.Context, key string)
}

// RateLimiter manages rate limiting for requests
type RateLimiter struct {
	config   *RateLimitConfig
	limiters sync.Map // map[string]*rate.Limiter for in-memory limiting
	redisLua *redis.Script
	ctx      context.Context
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	rl := &RateLimiter{
		config: config,
		ctx:    context.Background(),
	}

	// Initialize Redis Lua script for atomic operations
	if config.Redis != nil {
		rl.redisLua = redis.NewScript(`
			local key = KEYS[1]
			local limit = tonumber(ARGV[1])
			local window = tonumber(ARGV[2])
			local current = redis.call('INCR', key)
			
			if current == 1 then
				redis.call('EXPIRE', key, window)
			end
			
			if current > limit then
				return 0
			end
			
			return limit - current + 1
		`)
	}

	return rl
}

// Middleware returns the Echo middleware function
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !rl.config.Enabled {
				return next(c)
			}

			key := rl.getKey(c)

			allowed, remaining, resetTime := rl.allowRequest(key)

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerSecond))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

			if !allowed {
				// Call custom callback if provided
				if rl.config.OnLimitReached != nil {
					rl.config.OnLimitReached(c, key)
				}

				rl.config.Logger.Warn("rate limit exceeded",
					zap.String("key", key),
					zap.String("ip", c.RealIP()),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
				)

				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"error":       "Rate limit exceeded",
					"retry_after": resetTime - time.Now().Unix(),
					"limit":       rl.config.RequestsPerSecond,
				})
			}

			return next(c)
		}
	}
}

// getKey generates the rate limit key based on strategy
func (rl *RateLimiter) getKey(c echo.Context) string {
	if rl.config.CustomKeyFunc != nil {
		return rl.config.CustomKeyFunc(c)
	}

	switch rl.config.Strategy {
	case StrategyUser:
		userID := c.Get("userID")
		if userID != nil {
			return fmt.Sprintf("ratelimit:user:%v", userID)
		}
		// Fallback to IP if user not authenticated
		return fmt.Sprintf("ratelimit:ip:%s", c.RealIP())

	case StrategyIPAndUser:
		userID := c.Get("userID")
		if userID != nil {
			return fmt.Sprintf("ratelimit:ip_user:%s:%v", c.RealIP(), userID)
		}
		return fmt.Sprintf("ratelimit:ip:%s", c.RealIP())

	case StrategyEndpoint:
		return fmt.Sprintf("ratelimit:endpoint:%s:%s", c.Request().Method, c.Path())

	case StrategyIP:
		fallthrough
	default:
		return fmt.Sprintf("ratelimit:ip:%s", c.RealIP())
	}
}

// allowRequest checks if request is allowed based on rate limit
func (rl *RateLimiter) allowRequest(key string) (allowed bool, remaining int, resetTime int64) {
	if rl.config.Redis != nil {
		return rl.allowRequestRedis(key)
	}
	return rl.allowRequestInMemory(key)
}

// allowRequestRedis implements Redis-based distributed rate limiting
func (rl *RateLimiter) allowRequestRedis(key string) (bool, int, int64) {
	window := 1 // 1 second window

	result, err := rl.redisLua.Run(
		rl.ctx,
		rl.config.Redis,
		[]string{key},
		rl.config.RequestsPerSecond,
		window,
	).Int()

	if err != nil {
		rl.config.Logger.Error("redis rate limit error", zap.Error(err))
		// Fail open - allow request if Redis is down
		return true, rl.config.RequestsPerSecond, time.Now().Unix() + 1
	}

	if result == 0 {
		// Rate limit exceeded
		ttl, _ := rl.config.Redis.TTL(rl.ctx, key).Result()
		resetTime := time.Now().Add(ttl).Unix()
		return false, 0, resetTime
	}

	// Request allowed
	remaining := result
	resetTime := time.Now().Add(time.Second).Unix()
	return true, remaining, resetTime
}

// allowRequestInMemory implements in-memory token bucket rate limiting
func (rl *RateLimiter) allowRequestInMemory(key string) (bool, int, int64) {
	limiter := rl.getLimiter(key)

	reservation := limiter.Reserve()
	if !reservation.OK() {
		return false, 0, time.Now().Unix() + 1
	}

	delay := reservation.Delay()
	if delay > 0 {
		reservation.Cancel()
		resetTime := time.Now().Add(delay).Unix()
		return false, 0, resetTime
	}

	// Calculate remaining tokens
	remaining := int(limiter.Burst()) - int(limiter.Tokens())
	if remaining < 0 {
		remaining = 0
	}

	resetTime := time.Now().Add(time.Second).Unix()
	return true, remaining, resetTime
}

// getLimiter retrieves or creates a limiter for the given key
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	if limiter, exists := rl.limiters.Load(key); exists {
		return limiter.(*rate.Limiter)
	}

	// Create new limiter
	limiter := rate.NewLimiter(
		rate.Limit(rl.config.RequestsPerSecond),
		rl.config.Burst,
	)

	// Store and return
	actual, _ := rl.limiters.LoadOrStore(key, limiter)
	return actual.(*rate.Limiter)
}

// CleanupInMemoryLimiters removes old in-memory limiters (call periodically)
func (rl *RateLimiter) CleanupInMemoryLimiters() {
	if rl.config.Redis != nil {
		// No cleanup needed for Redis-based limiting
		return
	}

	// Remove limiters that haven't been used recently
	rl.limiters.Range(func(key, value interface{}) bool {
		limiter := value.(*rate.Limiter)
		if limiter.Tokens() >= float64(limiter.Burst()) {
			// Limiter is full (unused), safe to remove
			rl.limiters.Delete(key)
		}
		return true
	})
}

// RateLimitByEndpoint creates a rate limiter for specific endpoints
func RateLimitByEndpoint(redis *redis.Client, logger *zap.Logger, limits map[string]int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			endpoint := fmt.Sprintf("%s:%s", c.Request().Method, c.Path())

			limit, exists := limits[endpoint]
			if !exists {
				return next(c)
			}

			rl := NewRateLimiter(&RateLimitConfig{
				Enabled:           true,
				Strategy:          StrategyEndpoint,
				RequestsPerSecond: limit,
				Burst:             limit * 2,
				Redis:             redis,
				Logger:            logger,
			})

			return rl.Middleware()(next)(c)
		}
	}
}

// RateLimitByIP creates a simple IP-based rate limiter
func RateLimitByIP(requestsPerSecond, burst int, redis *redis.Client, logger *zap.Logger) echo.MiddlewareFunc {
	rl := NewRateLimiter(&RateLimitConfig{
		Enabled:           true,
		Strategy:          StrategyIP,
		RequestsPerSecond: requestsPerSecond,
		Burst:             burst,
		Redis:             redis,
		Logger:            logger,
	})

	return rl.Middleware()
}

// RateLimitByUser creates a user-based rate limiter
func RateLimitByUser(requestsPerSecond, burst int, redis *redis.Client, logger *zap.Logger) echo.MiddlewareFunc {
	rl := NewRateLimiter(&RateLimitConfig{
		Enabled:           true,
		Strategy:          StrategyUser,
		RequestsPerSecond: requestsPerSecond,
		Burst:             burst,
		Redis:             redis,
		Logger:            logger,
	})

	return rl.Middleware()
}
