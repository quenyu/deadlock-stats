package handlers

import (
	"net/http"

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
	steamID := c.Param("steamId")
	if steamID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Steam ID is required"})
	}

	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	if profile == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *PlayerProfileHandler) GetPlayerProfileV2(c echo.Context) error {
	steamID := c.Param("steamId")
	if steamID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Steam ID is required"})
	}

	profile, err := h.service.GetExtendedPlayerProfile(c.Request().Context(), steamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	if profile == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Player not found"})
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *PlayerProfileHandler) GetRecentMatches(c echo.Context) error {
	steamID := c.Param("steamId")
	if steamID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Steam ID is required"})
	}

	matches, err := h.service.GetRecentMatches(c.Request().Context(), steamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, matches)
}
