package security

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HeadersConfig struct {
	// HSTS (HTTP Strict Transport Security)
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	HSTSPreload           bool

	// XSS Protection
	XSSProtection      bool
	XFrameOptions      string // DENY, SAMEORIGIN
	ContentTypeNoSniff bool

	// Referrer Policy
	ReferrerPolicy string

	// Permissions Policy
	PermissionsPolicy string

	// Additional headers
	XContentTypeOptions   string
	XDNSPrefetchControl   string
	XDownloadOptions      string
	XPermittedCrossDomain string

	// Custom headers
	ServerHeader  string
	XPoweredBy    string
	RemoveHeaders []string

	Logger *zap.Logger
}

type HeadersMiddleware struct {
	config *HeadersConfig
	logger *zap.Logger
}

func NewHeadersMiddleware(config *HeadersConfig) *HeadersMiddleware {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	if config.XFrameOptions == "" {
		config.XFrameOptions = "DENY"
	}
	if config.ReferrerPolicy == "" {
		config.ReferrerPolicy = "strict-origin-when-cross-origin"
	}
	if config.HSTSMaxAge == 0 {
		config.HSTSMaxAge = 31536000 // 1 year
	}

	return &HeadersMiddleware{
		config: config,
		logger: config.Logger,
	}
}

func (h *HeadersMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h.setHeaders(c)
			return next(c)
		}
	}
}

func (h *HeadersMiddleware) setHeaders(c echo.Context) {
	res := c.Response()

	for _, header := range h.config.RemoveHeaders {
		res.Header().Del(header)
	}

	if h.config.ServerHeader != "" {
		res.Header().Set("Server", h.config.ServerHeader)
	} else {
		res.Header().Del("Server")
	}

	// X-Powered-By
	if h.config.XPoweredBy != "" {
		res.Header().Set("X-Powered-By", h.config.XPoweredBy)
	} else {
		res.Header().Del("X-Powered-By")
	}

	// HSTS (only for HTTPS)
	if h.isHTTPS(c) {
		hstsValue := fmt.Sprintf("max-age=%d", h.config.HSTSMaxAge)
		if h.config.HSTSIncludeSubdomains {
			hstsValue += "; includeSubDomains"
		}
		if h.config.HSTSPreload {
			hstsValue += "; preload"
		}
		res.Header().Set("Strict-Transport-Security", hstsValue)
	}

	// XSS Protection
	if h.config.XSSProtection {
		res.Header().Set("X-XSS-Protection", "1; mode=block")
	}

	// X-Frame-Options
	if h.config.XFrameOptions != "" {
		res.Header().Set("X-Frame-Options", h.config.XFrameOptions)
	}

	// X-Content-Type-Options
	if h.config.ContentTypeNoSniff {
		res.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// Referrer-Policy
	if h.config.ReferrerPolicy != "" {
		res.Header().Set("Referrer-Policy", h.config.ReferrerPolicy)
	}

	// Permissions-Policy
	if h.config.PermissionsPolicy != "" {
		res.Header().Set("Permissions-Policy", h.config.PermissionsPolicy)
	}

	// Additional security headers
	if h.config.XContentTypeOptions != "" {
		res.Header().Set("X-Content-Type-Options", h.config.XContentTypeOptions)
	}
	if h.config.XDNSPrefetchControl != "" {
		res.Header().Set("X-DNS-Prefetch-Control", h.config.XDNSPrefetchControl)
	}
	if h.config.XDownloadOptions != "" {
		res.Header().Set("X-Download-Options", h.config.XDownloadOptions)
	}
	if h.config.XPermittedCrossDomain != "" {
		res.Header().Set("X-Permitted-Cross-Domain-Policies", h.config.XPermittedCrossDomain)
	}
}

// isHTTPS checks if the connection is HTTPS
func (h *HeadersMiddleware) isHTTPS(c echo.Context) bool {
	return c.Request().TLS != nil || c.Request().Header.Get("X-Forwarded-Proto") == "https"
}
