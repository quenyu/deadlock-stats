package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/domain"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
	"github.com/quenyu/deadlock-stats/internal/services"
	"github.com/quenyu/deadlock-stats/internal/validators"
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

	matches, err := h.service.GetRecentMatches(c.Request().Context(), steamID)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if matches == nil {
		matches = []domain.Match{}
	}

	response := echo.Map{
		"matches": matches,
		"total":   len(matches),
		"page":    1,
		"limit":   limit,
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
	if err := validators.ValidateSteamID(steamID); err != nil {
		return "", cErrors.ErrInvalidSteamID
	}
	return steamID, nil
}

func (h *PlayerProfileHandler) validateSearchParams(c echo.Context) (string, string, error) {
	query := c.QueryParam("q")

	if err := validators.ValidatePlayerSearchQuery(query); err != nil {
		return "", "", err
	}

	searchType := c.QueryParam("type")
	if searchType == "" {
		searchType = "nickname"
	}

	return query, searchType, nil
}
