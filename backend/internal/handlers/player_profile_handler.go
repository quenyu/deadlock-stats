package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type PlayerProfileHandler struct {
	service *services.PlayerProfileService
}

func NewPlayerProfileHandler(service *services.PlayerProfileService) *PlayerProfileHandler {
	return &PlayerProfileHandler{
		service: service,
	}
}

func (h *PlayerProfileHandler) GetPlayerProfileV2(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return err
	}

	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return h.handleProfileResponse(c, profile)
}

func (h *PlayerProfileHandler) GetPlayerProfileWithMetrics(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return err
	}

	start := time.Now()
	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	if profile == nil {
		return h.handleServiceError(c, cErrors.ErrPlayerNotFound)
	}

	response := echo.Map{
		"steamID":  steamID,
		"loadTime": loadTime.Milliseconds(),
		"cacheHit": profile.LastUpdatedAt != time.Time{},
		"profile":  profile,
		"metrics": echo.Map{
			"totalLoadTime": loadTime.Milliseconds(),
			"hasData":       profile != nil,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerProfileHandler) GetRecentMatches(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return err
	}

	limitStr := c.QueryParam("limit")
	limit := 5 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	start := time.Now()
	matches, err := h.service.GetRecentMatches(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := echo.Map{
		"steamID":  steamID,
		"limit":    limit,
		"loadTime": loadTime.Milliseconds(),
		"matches":  matches,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerProfileHandler) SearchPlayers(c echo.Context) error {
	query, searchType, err := h.validateSearchParams(c)
	if err != nil {
		return err
	}

	users, err := h.service.SearchPlayers(c.Request().Context(), query, searchType)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *PlayerProfileHandler) validateSteamIDParam(c echo.Context) (string, error) {
	steamID := c.Param("steamId")
	if steamID == "" {
		return "", cErrors.ErrInvalidSteamID
	}
	return steamID, nil
}

func (h *PlayerProfileHandler) validateSearchParams(c echo.Context) (string, string, error) {
	query := c.QueryParam("q")
	if query == "" {
		return "", "", cErrors.ErrInvalidQuery
	}

	searchType := c.QueryParam("type")
	if searchType == "" {
		searchType = "nickname"
	}

	return query, searchType, nil
}

func (h *PlayerProfileHandler) handleProfileResponse(c echo.Context, profile interface{}) error {
	if profile == nil {
		return h.handleServiceError(c, cErrors.ErrPlayerNotFound)
	}
	return c.JSON(http.StatusOK, profile)
}

func (h *PlayerProfileHandler) handleServiceError(c echo.Context, err error) error {
	switch {
	case err == cErrors.ErrPlayerNotFound:
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
	case err == cErrors.ErrInvalidSteamID:
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Steam ID"})
	case err == cErrors.ErrInvalidQuery:
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid query parameter"})
	case err == cErrors.ErrRateLimited:
		return c.JSON(http.StatusTooManyRequests, echo.Map{"error": "Rate limit exceeded"})
	default:
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
}
