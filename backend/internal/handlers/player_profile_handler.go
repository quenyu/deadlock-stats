package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type PlayerProfileHandler struct {
	service *services.PlayerProfileService
}

func NewPlayerProfileHandler(service *services.PlayerProfileService) *PlayerProfileHandler {
	return &PlayerProfileHandler{service: service}
}

func (h *PlayerProfileHandler) GetPlayerProfile(c echo.Context) error {
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

func (h *PlayerProfileHandler) GetRecentMatches(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return err
	}

	matches, err := h.service.GetRecentMatches(c.Request().Context(), steamID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, matches)
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
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
	}

	response := echo.Map{
		"profile": profile,
		"metrics": echo.Map{
			"load_time_ms": loadTime.Milliseconds(),
			"cache_hit":    profile.LastUpdatedAt.IsZero(),
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerProfileHandler) validateSteamIDParam(c echo.Context) (string, error) {
	steamID := c.Param("steamId")
	if steamID == "" {
		return "", c.JSON(http.StatusBadRequest, echo.Map{"error": "Steam ID is required"})
	}
	return steamID, nil
}

func (h *PlayerProfileHandler) validateSearchParams(c echo.Context) (string, string, error) {
	query := c.QueryParam("q")
	if query == "" {
		return "", "", c.JSON(http.StatusBadRequest, echo.Map{"error": "Query is required"})
	}

	searchType := c.QueryParam("type")
	if searchType == "" {
		searchType = "nickname"
	}

	return query, searchType, nil
}

func (h *PlayerProfileHandler) handleServiceError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
}

func (h *PlayerProfileHandler) handleProfileResponse(c echo.Context, profile interface{}) error {
	if profile == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
	}
	return c.JSON(http.StatusOK, profile)
}
