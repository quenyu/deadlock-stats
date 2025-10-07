package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
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
		return ErrorHandler(cErrors.ErrInvalidRequestBody, c)
	}
	authorID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidUserID, c)
	}
	crosshair, err := h.service.Create(authorID, &req)
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusCreated, crosshair)
}

func (h *CrosshairHandler) GetAll(c echo.Context) error {
	page := 1
	limit := 20
	if p := c.QueryParam("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil {
			return ErrorHandler(cErrors.ErrInvalidRequestBody, c)
		}
	}
	if l := c.QueryParam("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			return ErrorHandler(cErrors.ErrInvalidRequestBody, c)
		}
	}
	crosshairs, err := h.service.GetAll(page, limit)
	if err != nil {
		return ErrorHandler(err, c)
	}

	total, err := h.service.Count()
	if err != nil {
		return ErrorHandler(err, c)
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
		return ErrorHandler(cErrors.ErrInvalidCrosshairID, c)
	}
	crosshair, err := h.service.GetByID(id)
	if err != nil {
		return ErrorHandler(cErrors.ErrCrosshairNotFound, c)
	}
	return c.JSON(http.StatusOK, crosshair)
}

func (h *CrosshairHandler) Like(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidCrosshairID, c)
	}
	userID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidUserID, c)
	}
	if err := h.service.Like(id, userID); err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Liked successfully"})
}

func (h *CrosshairHandler) Unlike(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidCrosshairID, c)
	}
	userID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidUserID, c)
	}
	if err := h.service.Unlike(id, userID); err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Unliked successfully"})
}

func (h *CrosshairHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidCrosshairID, c)
	}
	authorID, err := uuid.Parse(c.Get("userID").(string))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidUserID, c)
	}
	if err := h.service.Delete(id, authorID); err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Deleted successfully"})
}

func (h *CrosshairHandler) GetByAuthorID(c echo.Context) error {
	authorID, err := uuid.Parse(c.Param("author_id"))
	if err != nil {
		return ErrorHandler(cErrors.ErrInvalidUserID, c)
	}
	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			return ErrorHandler(cErrors.ErrInvalidQuery, c)
		}
	}
	crosshairs, err := h.service.GetByAuthorID(authorID, limit)
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, crosshairs)
}
