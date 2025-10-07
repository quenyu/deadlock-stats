package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/services"
)

type CrosshairHandler struct {
	service *services.CrosshairService
}

func NewCrosshairHandler(service *services.CrosshairService) *CrosshairHandler {
	return &CrosshairHandler{service: service}
}

func (h *CrosshairHandler) Create(c echo.Context) error {
	var req services.CreateCrosshairRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	authorID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}
	crosshair, err := h.service.Create(authorID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, crosshair)
}

func (h *CrosshairHandler) GetAll(c echo.Context) error {
	page := 1
	limit := 20
	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	crosshairs, err := h.service.GetAll(page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	total, err := h.service.Count()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"crosshairs": crosshairs,
		"total":      total,
		"page":       page,
		"limit":      limit,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *CrosshairHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crosshair ID")
	}
	crosshair, err := h.service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Crosshair not found")
	}
	return c.JSON(http.StatusOK, crosshair)
}

func (h *CrosshairHandler) Like(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crosshair ID")
	}
	userID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}
	if err := h.service.Like(id, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Liked successfully"})
}

func (h *CrosshairHandler) Unlike(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crosshair ID")
	}
	userID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}
	if err := h.service.Unlike(id, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Unliked successfully"})
}

func (h *CrosshairHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid crosshair ID")
	}
	authorID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user id")
	}
	if err := h.service.Delete(id, authorID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Deleted successfully"})
}

func (h *CrosshairHandler) GetByAuthorID(c echo.Context) error {
	authorID, err := uuid.Parse(c.Param("author_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid author ID")
	}
	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	crosshairs, err := h.service.GetByAuthorID(authorID, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, crosshairs)
}
