package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (s *AuthHandler) LoginHandler(c echo.Context) error {
	steamAuthURL, err := s.authService.InitiateSteamAuth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initiate Steam authentication"})
	}
	return c.Redirect(http.StatusTemporaryRedirect, steamAuthURL)
}

func (s *AuthHandler) CallbackHandler(c echo.Context) error {
	_, err := s.authService.HandleSteamCallback(c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to handle Steam callback"})
	}

	// TODO: Get JWT token of the user
	// TODO: Set JWT token in cookie
	// cookie := new(http.Cookie)
	// cookie.Name = "jwt"
	// cookie.Value = jwtToken
	// cookie.Expires = ...
	// c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
