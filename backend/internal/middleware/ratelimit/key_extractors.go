package ratelimit

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func IPKeyExtractor(c echo.Context) string {
	return c.RealIP()
}

// UserKeyExtractor extracts user ID as key (from JWT or context)
func UserKeyExtractor(c echo.Context) string {
	// Try to get user ID from context (set by JWT middleware)
	if userID := c.Get("user_id"); userID != nil {
		return fmt.Sprintf("user:%v", userID)
	}

	// Fallback to IP if user is not authenticated
	return IPKeyExtractor(c)
}

// IPAndUserKeyExtractor combines IP and user ID
func IPAndUserKeyExtractor(c echo.Context) string {
	ip := c.RealIP()

	if userID := c.Get("user_id"); userID != nil {
		return fmt.Sprintf("ip:%s:user:%v", ip, userID)
	}

	return fmt.Sprintf("ip:%s", ip)
}

// EndpointKeyExtractor extracts endpoint (method + path) as key
func EndpointKeyExtractor(c echo.Context) string {
	method := c.Request().Method
	path := c.Path()
	return fmt.Sprintf("%s:%s", method, path)
}

// IPAndEndpointKeyExtractor combines IP and endpoint
func IPAndEndpointKeyExtractor(c echo.Context) string {
	ip := c.RealIP()
	method := c.Request().Method
	path := c.Path()
	return fmt.Sprintf("ip:%s:%s:%s", ip, method, path)
}

// UserAndEndpointKeyExtractor combines user and endpoint
func UserAndEndpointKeyExtractor(c echo.Context) string {
	method := c.Request().Method
	path := c.Path()

	if userID := c.Get("user_id"); userID != nil {
		return fmt.Sprintf("user:%v:%s:%s", userID, method, path)
	}

	return fmt.Sprintf("ip:%s:%s:%s", c.RealIP(), method, path)
}

// GetKeyExtractor returns key extractor based on strategy
func GetKeyExtractor(strategy Strategy) KeyExtractor {
	switch strategy {
	case StrategyIP:
		return IPKeyExtractor
	case StrategyUser:
		return UserKeyExtractor
	case StrategyIPAndUser:
		return IPAndUserKeyExtractor
	case StrategyEndpoint:
		return EndpointKeyExtractor
	default:
		return IPKeyExtractor
	}
}

// GetEndpointKey returns normalized endpoint key for per-endpoint limits
func GetEndpointKey(c echo.Context) string {
	method := c.Request().Method
	path := c.Path()

	// Normalize path (remove trailing slash)
	path = strings.TrimSuffix(path, "/")

	return fmt.Sprintf("%s:%s", method, path)
}
