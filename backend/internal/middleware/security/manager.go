package security

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ManagerConfig struct {
	Headers *HeadersConfig
	CSP     *CSPConfig
	CORS    *CORSConfig
	CSRF    *CSRFConfig
	Logger  *zap.Logger
}

type Manager struct {
	config  *ManagerConfig
	headers *HeadersMiddleware
	csp     *CSPMiddleware
	cors    *CORSMiddleware
	csrf    *CSRFMiddleware
	logger  *zap.Logger
}

func NewManager(config *ManagerConfig) *Manager {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	manager := &Manager{
		config: config,
		logger: config.Logger,
	}

	if config.Headers != nil {
		manager.headers = NewHeadersMiddleware(config.Headers)
	}
	if config.CSP != nil {
		manager.csp = NewCSPMiddleware(config.CSP)
	}
	if config.CORS != nil {
		manager.cors = NewCORSMiddleware(config.CORS)
	}
	if config.CSRF != nil {
		manager.csrf = NewCSRFMiddleware(config.CSRF)
	}

	return manager
}

// Middleware returns the combined security middleware
func (m *Manager) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := next

			// CSRF (last - validates before processing)
			if m.csrf != nil {
				handler = m.csrf.Middleware()(handler)
			}

			// CORS (before CSRF to handle preflight)
			if m.cors != nil {
				handler = m.cors.Middleware()(handler)
			}

			// CSP
			if m.csp != nil {
				handler = m.csp.Middleware()(handler)
			}

			// Headers (first - sets base security headers)
			if m.headers != nil {
				handler = m.headers.Middleware()(handler)
			}

			return handler(c)
		}
	}
}

// GetCSRFMiddleware returns the CSRF middleware for manual token generation
func (m *Manager) GetCSRFMiddleware() *CSRFMiddleware {
	return m.csrf
}

// DefaultConfig returns production-ready security configuration
func DefaultConfig(clientURL string, logger *zap.Logger) *ManagerConfig {
	return &ManagerConfig{
		Headers: &HeadersConfig{
			HSTSMaxAge:            31536000, // 1 year
			HSTSIncludeSubdomains: true,
			HSTSPreload:           false,
			XSSProtection:         true,
			XFrameOptions:         "DENY",
			ContentTypeNoSniff:    true,
			ReferrerPolicy:        "strict-origin-when-cross-origin",
			PermissionsPolicy:     "geolocation=(), microphone=(), camera=()",
			XContentTypeOptions:   "nosniff",
			XDNSPrefetchControl:   "off",
			XDownloadOptions:      "noopen",
			XPermittedCrossDomain: "none",
			RemoveHeaders:         []string{"Server", "X-Powered-By"},
			Logger:                logger,
		},
		CSP: &CSPConfig{
			Enabled:    true,
			ReportOnly: false,
			Directives: DefaultCSPDirectives(),
			Logger:     logger,
		},
		CORS: &CORSConfig{
			Enabled:          true,
			AllowOrigins:     []string{clientURL},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
			ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
			AllowCredentials: true,
			MaxAge:           86400,
			Logger:           logger,
		},
		CSRF: &CSRFConfig{
			Enabled:        true,
			TokenLength:    32,
			CookieName:     "_csrf",
			HeaderName:     "X-CSRF-Token",
			CookiePath:     "/",
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: http.SameSiteStrictMode,
			TokenLookup:    "header:X-CSRF-Token",
			SkipPaths:      []string{"/health", "/metrics"},
			Logger:         logger,
		},
		Logger: logger,
	}
}

// DevelopmentConfig returns configuration suitable for development
func DevelopmentConfig(clientURL string, logger *zap.Logger) *ManagerConfig {
	config := DefaultConfig(clientURL, logger)

	// Relax HTTPS requirements
	config.Headers.HSTSMaxAge = 0 // Disable HSTS
	config.CSRF.CookieSecure = false
	config.CSP.ReportOnly = true // Report-only mode

	return config
}
