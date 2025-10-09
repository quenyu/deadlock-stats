package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CSRFConfig struct {
	Enabled        bool
	TokenLength    int
	CookieName     string
	HeaderName     string
	CookiePath     string
	CookieDomain   string
	CookieSecure   bool
	CookieHTTPOnly bool
	CookieSameSite http.SameSite
	TokenLookup    string // "header:X-CSRF-Token" or "form:csrf"
	SkipPaths      []string
	Logger         *zap.Logger
}

type CSRFMiddleware struct {
	config *CSRFConfig
	tokens sync.Map // map[string]time.Time
	logger *zap.Logger
}

func NewCSRFMiddleware(config *CSRFConfig) *CSRFMiddleware {
	if config.Logger == nil {
		config.Logger, _ = zap.NewProduction()
	}

	if config.TokenLength == 0 {
		config.TokenLength = 32
	}
	if config.CookieName == "" {
		config.CookieName = "_csrf"
	}
	if config.HeaderName == "" {
		config.HeaderName = "X-CSRF-Token"
	}
	if config.CookiePath == "" {
		config.CookiePath = "/"
	}
	if config.TokenLookup == "" {
		config.TokenLookup = "header:X-CSRF-Token"
	}

	csrf := &CSRFMiddleware{
		config: config,
		logger: config.Logger,
	}

	go csrf.cleanupExpiredTokens()

	return csrf
}

// Middleware returns the Echo middleware function
func (csrf *CSRFMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !csrf.config.Enabled {
				return next(c)
			}

			if csrf.isSafeMethod(c.Request().Method) {
				return next(c)
			}

			if csrf.isPathSkipped(c.Path()) {
				return next(c)
			}

			if err := csrf.verifyToken(c); err != nil {
				csrf.logger.Warn("CSRF verification failed",
					zap.String("ip", c.RealIP()),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
					zap.Error(err),
				)
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token verification failed")
			}

			return next(c)
		}
	}
}

// GenerateToken generates and sets a new CSRF token
func (csrf *CSRFMiddleware) GenerateToken(c echo.Context) (string, error) {
	token, err := generateRandomToken(csrf.config.TokenLength)
	if err != nil {
		return "", err
	}

	// Store token with expiration (24 hours)
	csrf.tokens.Store(token, time.Now().Add(24*time.Hour))

	cookie := &http.Cookie{
		Name:     csrf.config.CookieName,
		Value:    token,
		Path:     csrf.config.CookiePath,
		Domain:   csrf.config.CookieDomain,
		MaxAge:   86400, // 24 hours
		Secure:   csrf.config.CookieSecure,
		HttpOnly: csrf.config.CookieHTTPOnly,
		SameSite: csrf.config.CookieSameSite,
	}

	c.SetCookie(cookie)

	return token, nil
}

// verifyToken verifies the CSRF token
func (csrf *CSRFMiddleware) verifyToken(c echo.Context) error {
	cookie, err := c.Cookie(csrf.config.CookieName)
	if err != nil {
		return fmt.Errorf("CSRF cookie not found")
	}

	cookieToken := cookie.Value

	requestToken := csrf.extractTokenFromRequest(c)
	if requestToken == "" {
		return fmt.Errorf("CSRF token not provided")
	}
	if cookieToken != requestToken {
		return fmt.Errorf("CSRF token mismatch")
	}

	if _, exists := csrf.tokens.Load(cookieToken); !exists {
		return fmt.Errorf("CSRF token expired or invalid")
	}

	return nil
}

// extractTokenFromRequest extracts token from request based on TokenLookup
func (csrf *CSRFMiddleware) extractTokenFromRequest(c echo.Context) string {
	if strings.HasPrefix(csrf.config.TokenLookup, "header:") {
		headerName := strings.TrimPrefix(csrf.config.TokenLookup, "header:")
		return c.Request().Header.Get(headerName)
	} else if strings.HasPrefix(csrf.config.TokenLookup, "form:") {
		formField := strings.TrimPrefix(csrf.config.TokenLookup, "form:")
		return c.FormValue(formField)
	}

	return c.Request().Header.Get(csrf.config.HeaderName)
}

// isSafeMethod checks if HTTP method is safe (doesn't need CSRF)
func (csrf *CSRFMiddleware) isSafeMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions ||
		method == http.MethodTrace
}

// isPathSkipped checks if path should skip CSRF verification
func (csrf *CSRFMiddleware) isPathSkipped(path string) bool {
	for _, skipPath := range csrf.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// cleanupExpiredTokens removes expired tokens periodically
func (csrf *CSRFMiddleware) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		csrf.tokens.Range(func(key, value interface{}) bool {
			expiresAt := value.(time.Time)
			if now.After(expiresAt) {
				csrf.tokens.Delete(key)
			}
			return true
		})
	}
}

// generateRandomToken generates cryptographically secure random token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
