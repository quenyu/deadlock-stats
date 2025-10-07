package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
)

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var errorMap = map[error]httpError{
	// Player-related
	cErrors.ErrPlayerNotFound: {http.StatusNotFound, "Player not found"},
	cErrors.ErrInvalidSteamID: {http.StatusBadRequest, "Invalid Steam ID"},
	cErrors.ErrRateLimited:    {http.StatusTooManyRequests, "Rate limited"},
	cErrors.ErrAPIUnavailable: {http.StatusServiceUnavailable, "External API unavailable"},

	// Auth-related
	cErrors.ErrUserNotFound:        {http.StatusNotFound, "User not found"},
	cErrors.ErrInvalidCredentials:  {http.StatusBadRequest, "Invalid credentials"},
	cErrors.ErrUnauthorized:        {http.StatusUnauthorized, "Unauthorized"},
	cErrors.ErrForbidden:           {http.StatusForbidden, "Forbidden"},
	cErrors.ErrInvalidToken:        {http.StatusUnauthorized, "Invalid token"},
	cErrors.ErrSessionExpired:      {http.StatusUnauthorized, "Session expired"},
	cErrors.ErrSteamAuthFailed:     {http.StatusUnauthorized, "Steam authentication failed"},
	cErrors.ErrJWTGenerationFailed: {http.StatusInternalServerError, "Failed to generate JWT token"},
}

func ErrorHandler(err error, c echo.Context) error {
	for targetErr, httpErr := range errorMap {
		if errors.Is(err, targetErr) {
			return c.JSON(httpErr.Code, echo.Map{
				"error": httpErr.Message,
				"code":  httpErr.Code,
			})
		}
	}

	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error": "Internal server error",
		"code":  http.StatusInternalServerError,
	})
}
