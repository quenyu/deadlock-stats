package security

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CORSConfig struct {
	Enabled          bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
	Logger           *zap.Logger
}

type CORSMiddleware struct {
	config *CORSConfig
	logger *zap.Logger
}

func NewCORSMiddleware(config *CORSConfig) *CORSMiddleware {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	if config.MaxAge == 0 {
		config.MaxAge = 86400 // 24 hours
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	return &CORSMiddleware{
		config: config,
		logger: config.Logger,
	}
}

// Middleware returns the Echo middleware function
func (cors *CORSMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !cors.config.Enabled {
				return next(c)
			}

			if c.Request().Method == http.MethodOptions {
				return cors.handlePreflight(c)
			}

			cors.setHeaders(c)

			return next(c)
		}
	}
}

// setHeaders sets CORS headers for regular requests
func (cors *CORSMiddleware) setHeaders(c echo.Context) {
	origin := c.Request().Header.Get("Origin")
	res := c.Response()

	if cors.isOriginAllowed(origin) {
		res.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(cors.config.AllowOrigins) > 0 && cors.config.AllowOrigins[0] == "*" {
		res.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Credentials
	if cors.config.AllowCredentials {
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Expose headers
	if len(cors.config.ExposeHeaders) > 0 {
		res.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.config.ExposeHeaders, ", "))
	}

	// Vary header for caching
	res.Header().Add("Vary", "Origin")
}

// handlePreflight handles OPTIONS preflight requests
func (cors *CORSMiddleware) handlePreflight(c echo.Context) error {
	origin := c.Request().Header.Get("Origin")
	res := c.Response()

	if !cors.isOriginAllowed(origin) {
		cors.logger.Warn("CORS preflight rejected",
			zap.String("origin", origin),
			zap.String("ip", c.RealIP()),
		)
		return echo.NewHTTPError(http.StatusForbidden, "Origin not allowed")
	}

	res.Header().Set("Access-Control-Allow-Origin", origin)

	if len(cors.config.AllowMethods) > 0 {
		res.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.config.AllowMethods, ", "))
	}

	if len(cors.config.AllowHeaders) > 0 {
		res.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.config.AllowHeaders, ", "))
	}

	if cors.config.AllowCredentials {
		res.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if cors.config.MaxAge > 0 {
		res.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cors.config.MaxAge))
	}

	return c.NoContent(http.StatusNoContent)
}

// isOriginAllowed checks if origin is in allowed list
func (cors *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowed := range cors.config.AllowOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// Support wildcard subdomains (*.example.com)
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}
