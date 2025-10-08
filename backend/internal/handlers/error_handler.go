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
	cErrors.ErrPlayerNotFound:    {http.StatusNotFound, "Player not found"},
	cErrors.ErrInvalidSteamID:    {http.StatusBadRequest, "Invalid Steam ID"},
	cErrors.ErrInvalidQuery:      {http.StatusBadRequest, "Invalid query parameter"},
	cErrors.ErrRateLimited:       {http.StatusTooManyRequests, "Rate limited"},
	cErrors.ErrAPIUnavailable:    {http.StatusServiceUnavailable, "External API unavailable"},
	cErrors.ErrPlayerDataMissing: {http.StatusNotFound, "Player data missing"},

	// Auth-related
	cErrors.ErrUserNotFound:        {http.StatusNotFound, "User not found"},
	cErrors.ErrInvalidCredentials:  {http.StatusBadRequest, "Invalid credentials"},
	cErrors.ErrUnauthorized:        {http.StatusUnauthorized, "Unauthorized"},
	cErrors.ErrForbidden:           {http.StatusForbidden, "Forbidden"},
	cErrors.ErrInvalidToken:        {http.StatusUnauthorized, "Invalid token"},
	cErrors.ErrSessionExpired:      {http.StatusUnauthorized, "Session expired"},
	cErrors.ErrSteamAuthFailed:     {http.StatusUnauthorized, "Steam authentication failed"},
	cErrors.ErrJWTGenerationFailed: {http.StatusInternalServerError, "Failed to generate JWT token"},

	// Crosshair-related
	cErrors.ErrCrosshairNotFound:  {http.StatusNotFound, "Crosshair not found"},
	cErrors.ErrInvalidCrosshairID: {http.StatusBadRequest, "Invalid crosshair ID"},
	cErrors.ErrInvalidUserID:      {http.StatusBadRequest, "Invalid user ID"},
	cErrors.ErrCrosshairForbidden: {http.StatusForbidden, "Forbidden"},
	cErrors.ErrInvalidRequestBody: {http.StatusBadRequest, "Invalid request body"},

	// Match-related
	cErrors.ErrMatchNotFound:   {http.StatusNotFound, "Match not found"},
	cErrors.ErrInvalidSearch:   {http.StatusBadRequest, "Invalid search type"},
	cErrors.ErrNoSearchResults: {http.StatusNotFound, "No search results"},

	// System-related
	cErrors.ErrDatabaseError:   {http.StatusInternalServerError, "Database operation failed"},
	cErrors.ErrCacheError:      {http.StatusInternalServerError, "Cache operation failed"},
	cErrors.ErrUnknownInternal: {http.StatusInternalServerError, "Unknown internal error"},
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
