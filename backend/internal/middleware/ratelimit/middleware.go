package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Middleware struct {
	config    *Config
	limiter   Limiter
	logger    *zap.Logger
	onLimited LimitReachedCallback
}

func NewMiddleware(config *Config, limiter Limiter, logger *zap.Logger) (*Middleware, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &Middleware{
		config:  config,
		limiter: limiter,
		logger:  logger,
	}, nil
}

func (m *Middleware) SetOnLimitReached(callback LimitReachedCallback) {
	m.onLimited = callback
}

// Handler returns the Echo middleware handler
func (m *Middleware) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip if disabled
			if !m.config.Enabled {
				return next(c)
			}

			// Check whitelist
			ip := c.RealIP()
			if m.config.IsWhitelisted(ip) {
				return next(c)
			}

			// Extract key using strategy
			var key string
			if m.config.KeyFunc != nil {
				key = m.config.KeyFunc(c)
			} else {
				keyExtractor := GetKeyExtractor(m.config.Strategy)
				key = keyExtractor(c)
			}

			// Determine limit for this request
			limit := m.config.RequestsPerSecond

			// Check for per-endpoint limit
			endpointKey := GetEndpointKey(c)
			if endpointLimit, ok := m.config.PerEndpoint[endpointKey]; ok {
				limit = endpointLimit
			}

			// Check rate limit
			allowed, remaining, resetAt, err := m.limiter.Allow(key, limit)
			if err != nil {
				m.logger.Error("rate limit check failed",
					zap.Error(err),
					zap.String("key", key),
				)
				// Fail open: allow request if rate limiter fails
				return next(c)
			}

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))

			// Check if allowed
			if !allowed {
				// Call callback if set
				if m.onLimited != nil {
					m.onLimited(c, key)
				}

				m.logger.Warn("rate limit exceeded",
					zap.String("key", key),
					zap.String("ip", ip),
					zap.String("endpoint", endpointKey),
					zap.Int("limit", limit),
				)

				c.Response().Header().Set("Retry-After", strconv.FormatInt(resetAt, 10))

				return echo.NewHTTPError(
					http.StatusTooManyRequests,
					map[string]interface{}{
						"error":       "rate limit exceeded",
						"limit":       limit,
						"reset_at":    resetAt,
						"retry_after": resetAt,
					},
				)
			}

			return next(c)
		}
	}
}

// Close closes the middleware and underlying limiter
func (m *Middleware) Close() error {
	if m.limiter != nil {
		return m.limiter.Close()
	}
	return nil
}
