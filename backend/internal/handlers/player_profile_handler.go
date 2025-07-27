package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	start := time.Now()
	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	if profile == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player profile not found"})
	}

	response := echo.Map{
		"steamID":  steamID,
		"loadTime": loadTime.Milliseconds(),
		"profile":  profile,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerProfileHandler) GetPlayerProfileWithMetrics(c echo.Context) error {
	steamID, err := h.validateSteamIDParam(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	start := time.Now()
	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	loadTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	if profile == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player profile not found"})
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
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

func (h *PlayerProfileHandler) validateSteamIDParam(c echo.Context) (string, error) {
	steamID := c.Param("steamId")
	if steamID == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "SteamID parameter is required")
	}
	return steamID, nil
}

func (h *PlayerProfileHandler) handleServiceError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error":   "Internal server error",
		"message": err.Error(),
	})
}
