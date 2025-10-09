package ratelimit

import (
	"github.com/labstack/echo/v4"
)

// KeyExtractor extracts rate limit key from request
type KeyExtractor func(c echo.Context) string

// LimitReachedCallback is called when rate limit is exceeded
type LimitReachedCallback func(c echo.Context, key string)

// Limiter interface defines rate limiting operations
type Limiter interface {
	// Allow checks if request is allowed
	Allow(key string, limit int) (allowed bool, remaining int, resetAt int64, err error)

	// Close closes the limiter
	Close() error
}
