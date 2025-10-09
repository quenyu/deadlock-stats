package security

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CSPConfig struct {
	Enabled    bool
	ReportOnly bool
	Directives map[string]string
	ReportURI  string
	Logger     *zap.Logger
}

type CSPMiddleware struct {
	config *CSPConfig
	logger *zap.Logger
}

func NewCSPMiddleware(config *CSPConfig) *CSPMiddleware {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	if len(config.Directives) == 0 {
		config.Directives = DefaultCSPDirectives()
	}

	return &CSPMiddleware{
		config: config,
		logger: config.Logger,
	}
}

// Middleware returns the Echo middleware function
func (csp *CSPMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if csp.config.Enabled {
				csp.setCSPHeader(c)
			}
			return next(c)
		}
	}
}

// setCSPHeader sets the Content-Security-Policy header
func (csp *CSPMiddleware) setCSPHeader(c echo.Context) {
	policy := csp.buildPolicy()
	headerName := "Content-Security-Policy"

	if csp.config.ReportOnly {
		headerName = "Content-Security-Policy-Report-Only"
	}

	c.Response().Header().Set(headerName, policy)
}

// buildPolicy constructs the CSP policy string
func (csp *CSPMiddleware) buildPolicy() string {
	var directives []string

	for key, value := range csp.config.Directives {
		directives = append(directives, fmt.Sprintf("%s %s", key, value))
	}

	policy := strings.Join(directives, "; ")

	if csp.config.ReportURI != "" {
		policy += fmt.Sprintf("; report-uri %s", csp.config.ReportURI)
	}

	return policy
}

// DefaultCSPDirectives returns secure default CSP directives
func DefaultCSPDirectives() map[string]string {
	return map[string]string{
		"default-src":     "'self'",
		"script-src":      "'self' 'unsafe-inline' 'unsafe-eval'",
		"style-src":       "'self' 'unsafe-inline'",
		"img-src":         "'self' data: https:",
		"font-src":        "'self' data:",
		"connect-src":     "'self'",
		"frame-ancestors": "'none'",
		"base-uri":        "'self'",
		"form-action":     "'self'",
		"object-src":      "'none'",
		"media-src":       "'self'",
		"worker-src":      "'self'",
		"manifest-src":    "'self'",
	}
}

// StrictCSPDirectives returns very strict CSP directives
func StrictCSPDirectives() map[string]string {
	return map[string]string{
		"default-src":     "'none'",
		"script-src":      "'self'",
		"style-src":       "'self'",
		"img-src":         "'self'",
		"font-src":        "'self'",
		"connect-src":     "'self'",
		"frame-ancestors": "'none'",
		"base-uri":        "'self'",
		"form-action":     "'self'",
		"object-src":      "'none'",
	}
}
