package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/config"
)

type JWTMiddleware struct {
	config *config.Config
}

func NewJWTMiddleware(cfg *config.Config) *JWTMiddleware {
	return &JWTMiddleware{config: cfg}
}

func (m *JWTMiddleware) Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr, err := m.extractToken(c)
		if err != nil {
			return m.unauthorizedResponse(c)
		}

		userID, err := m.validateToken(tokenStr)
		if err != nil {
			return m.unauthorizedResponse(c)
		}

		c.Set("userID", userID)
		return next(c)
	}
}

func (m *JWTMiddleware) extractToken(c echo.Context) (string, error) {
	tokenStr := c.Request().Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	if tokenStr != "" {
		return tokenStr, nil
	}

	if ck, err := c.Cookie("jwt"); err == nil {
		return ck.Value, nil
	}

	return "", echo.NewHTTPError(http.StatusUnauthorized, "no token provided")
}

func (m *JWTMiddleware) validateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(m.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID in token")
	}

	return userID, nil
}

func (m *JWTMiddleware) unauthorizedResponse(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
}
