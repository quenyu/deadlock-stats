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
		return ErrorHandler(err, c)
	}

	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	if err != nil {
		return ErrorHandler(err, c)
	}

	if profile == nil {
		return ErrorHandler(cErrors.ErrPlayerNotFound, c)
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *PlayerProfileHandler) GetPlayerProfileWithMetrics(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return ErrorHandler(err, c)
	}

	start := time.Now()
	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if profile == nil {
		return ErrorHandler(cErrors.ErrPlayerNotFound, c)
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
		return ErrorHandler(err, c)
	}

	limit := 5
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil || val < 1 || val > 20 {
			return ErrorHandler(cErrors.ErrInvalidQuery, c)
		}
		limit = val
	}

	start := time.Now()
	matches, err := h.service.GetRecentMatches(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if len(matches) == 0 {
		return ErrorHandler(cErrors.ErrMatchNotFound, c)
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
		return ErrorHandler(err, c)
	}

	users, err := h.service.SearchPlayers(c.Request().Context(), query, searchType)
	if err != nil {
		return ErrorHandler(err, c)
	}

	if len(users) == 0 {
		return ErrorHandler(cErrors.ErrNoSearchResults, c)
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
