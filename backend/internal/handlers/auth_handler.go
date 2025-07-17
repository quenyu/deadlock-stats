package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/config"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, config: config}
}

func (s *AuthHandler) LoginHandler(c echo.Context) error {
	steamAuthURL, err := s.authService.InitiateSteamAuth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initiate Steam authentication"})
	}
	return c.Redirect(http.StatusTemporaryRedirect, steamAuthURL)
}

func (s *AuthHandler) CallbackHandler(c echo.Context) error {
	jwtToken, err := s.authService.HandleSteamCallback(c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to handle Steam callback"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = jwtToken
	cookie.Expires = time.Now().Add(s.config.JWT.Expiration)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	// You can also set cookie.Domain, cookie.Secure in production
	c.SetCookie(cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
}

func (h *AuthHandler) GetMyProfileHandler(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user ID from context"})
	}

	// Here you would call a service to get user details by ID
	// For now, just return the ID
	return c.JSON(http.StatusOK, map[string]string{"message": "Authenticated successfully", "userID": userID})
}
