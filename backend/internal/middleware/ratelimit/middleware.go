package ratelimit

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Middleware represents the rate limiting middleware
type Middleware struct {
	config         *Config
	limiter        ILimiter
	logger         *zap.Logger
	onLimitReached LimitReachedCallback
	keyExtractor   KeyExtractor
}

// NewMiddleware creates a new rate limiting middleware
func NewMiddleware(config *Config, limiter ILimiter, logger *zap.Logger) (*Middleware, error) {
	if config == nil {
		return nil, ErrInvalidConfig
	}

	if limiter == nil {
		return nil, ErrInvalidConfig
	}

	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	// Set default key extractor based on strategy
	keyExtractor := config.KeyFunc
	if keyExtractor == nil {
		keyExtractor = GetKeyExtractor(config.Strategy)
	}

	return &Middleware{
		config:       config,
		limiter:      limiter,
		logger:       logger,
		keyExtractor: keyExtractor,
	}, nil
}

// Handler returns the Echo middleware function
func (m *Middleware) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip if rate limiting is disabled
			if !m.config.Enabled {
				return next(c)
			}

			// Extract key for rate limiting
			key := m.keyExtractor(c)
			if key == "" {
				m.logger.Warn("empty rate limit key, skipping rate limiting")
				return next(c)
			}

			// Check if IP is whitelisted
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			if m.config.IsWhitelisted(clientIP) {
				return next(c)
			}

			// Get limit for this endpoint
			endpoint := c.Request().URL.Path
			limit := m.config.GetLimit(endpoint)

			// Check rate limit
			allowed, remaining, resetAt, err := m.limiter.Allow(key, limit)
			if err != nil {
				m.logger.Error("rate limiter error",
					zap.String("key", key),
					zap.Error(err))
				// On error, allow the request but log it
				return next(c)
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))

			if !allowed {
				// Call limit reached callback if set
				if m.onLimitReached != nil {
					m.onLimitReached(c, key)
				}

				m.logger.Warn("rate limit exceeded",
					zap.String("key", key),
					zap.String("endpoint", endpoint),
					zap.Int("limit", limit),
					zap.Int("remaining", remaining))

				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "Rate limit exceeded",
					"code":        http.StatusTooManyRequests,
					"retry_after": int(resetAt - time.Now().Unix()),
					"limit":       limit,
					"remaining":   remaining,
				})
			}

			return next(c)
		}
	}
}

// SetOnLimitReached sets the callback for when rate limit is exceeded
func (m *Middleware) SetOnLimitReached(callback LimitReachedCallback) {
	m.onLimitReached = callback
}

// Close closes the middleware
func (m *Middleware) Close() error {
	return m.limiter.Close()
}

// RateLimitMiddleware creates middleware for rate limiting (legacy function)
func RateLimitMiddleware(rateLimit rate.Limit, burst int) echo.MiddlewareFunc {
	limiter := NewLimiter(rateLimit, burst)

	// Cleaning old limiters every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the client's IP address
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// Checking the limit using new interface
			allowed, _, _, err := limiter.Allow(clientIP, int(rateLimit))
			if err != nil {
				// On error, allow the request
				return next(c)
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error":       "Rate limit exceeded",
					"code":        http.StatusTooManyRequests,
					"retry_after": 60, // seconds
				})
			}

			return next(c)
		}
	}
}

// StrictRateLimitMiddleware Stricter limit for critical endpoints
func StrictRateLimitMiddleware() echo.MiddlewareFunc {
	// 10 запросов в минуту, burst 5
	return RateLimitMiddleware(rate.Every(6*time.Second), 5)
}

// StandardRateLimitMiddleware Standard limit for regular endpoints
func StandardRateLimitMiddleware() echo.MiddlewareFunc {
	// 100 RPM, burst 20
	return RateLimitMiddleware(rate.Every(600*time.Millisecond), 20)
}

// LenientRateLimitMiddleware Soft limit for public endpoints
func LenientRateLimitMiddleware() echo.MiddlewareFunc {
	// 300 RPM, burst 50
	return RateLimitMiddleware(rate.Every(200*time.Millisecond), 50)
}
