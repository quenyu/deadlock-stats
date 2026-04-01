package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/config"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, config: config}
}

func (h *AuthHandler) LoginHandler(c echo.Context) error {
	steamAuthURL, err := h.authService.InitiateSteamAuth()
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.Redirect(http.StatusTemporaryRedirect, steamAuthURL)
}

func (h *AuthHandler) CallbackHandler(c echo.Context) error {
	jwtToken, err := h.authService.HandleSteamCallback(c.Request())
	if err != nil {
		return ErrorHandler(err, c)
	}

	h.setJWTCookie(c, jwtToken)
	return c.Redirect(http.StatusTemporaryRedirect, h.config.App.ClientURL)
}

func (h *AuthHandler) LogoutHandler(c echo.Context) error {
	h.clearJWTCookie(c)
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}

func (h *AuthHandler) GetMyProfileHandler(c echo.Context) error {
	userID, err := h.extractUserIDFromContext(c)
	if err != nil {
		return ErrorHandler(err, c)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Authenticated successfully", "userID": userID})
}

func (h *AuthHandler) GetUserMe(c echo.Context) error {
	userID, err := h.extractUserIDFromContext(c)
	if err != nil {
		return ErrorHandler(err, c)
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		return ErrorHandler(err, c)
	}

	if user == nil {
		return ErrorHandler(cErrors.ErrUserNotFound, c)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) extractUserIDFromContext(c echo.Context) (string, error) {
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return "", cErrors.ErrInvalidToken
	}
	return userID, nil
}

func (h *AuthHandler) setJWTCookie(c echo.Context, token string) {
	cookie := h.createJWTCookie(token, h.config.JWT.Expiration)
	c.SetCookie(cookie)
}

func (h *AuthHandler) clearJWTCookie(c echo.Context) {
	cookie := h.createJWTCookie("", -1)
	c.SetCookie(cookie)
}

func (h *AuthHandler) createJWTCookie(value string, expiration time.Duration) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode

	if expiration > 0 {
		cookie.Expires = time.Now().Add(expiration)
	} else {
		cookie.Expires = time.Unix(0, 0)
		cookie.MaxAge = -1
	}

	return cookie
}
