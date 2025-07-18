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
		tokenStr := c.Request().Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		if tokenStr == "" {
			if ck, err := c.Cookie("jwt"); err == nil {
				tokenStr = ck.Value
			}
		}

		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.config.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		c.Set("userID", userID)
		return next(c)
	}
}
