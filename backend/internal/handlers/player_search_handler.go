package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/dto"
	"github.com/quenyu/deadlock-stats/internal/services"
	"go.uber.org/zap"
)

type PlayerSearchHandler struct {
	searchService *services.PlayerSearchService
	logger        *zap.Logger
}

func NewPlayerSearchHandler(searchService *services.PlayerSearchService, logger *zap.Logger) *PlayerSearchHandler {
	return &PlayerSearchHandler{
		searchService: searchService,
		logger:        logger,
	}
}

func (h *PlayerSearchHandler) SearchPlayers(c echo.Context) error {
	start := time.Now()
	query := c.QueryParam("query")
	searchType := c.QueryParam("searchType")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, err := h.searchService.SearchPlayers(c.Request().Context(), query, searchType, page, pageSize)
	searchTime := time.Since(start)
	if err != nil {
		h.logger.Error("SearchPlayers error", zap.Error(err))
		return c.JSON(500, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{
		"results":    result.Results,
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"searchType": searchType,
		"searchTime": searchTime.Milliseconds(),
	})
}

func (h *PlayerSearchHandler) SearchPlayersAutocomplete(c echo.Context) error {
	query := c.QueryParam("query")
	limitStr := c.QueryParam("limit")

	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Query parameter is required"})
	}

	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	start := time.Now()
	users, err := h.searchService.SearchPlayersWithAutocomplete(c.Request().Context(), query, limit)
	searchTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := echo.Map{
		"query":      query,
		"limit":      limit,
		"searchTime": searchTime.Milliseconds(),
		"totalFound": len(users),
		"results":    users,
	}

	// Логируем ответ для отладки
	h.logger.Info("Search autocomplete response",
		zap.String("query", query),
		zap.Int("totalFound", len(users)),
		zap.Any("results", users))

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) SearchPlayersWithFilters(c echo.Context) error {
	query := c.QueryParam("query")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Query parameter is required"})
	}

	filters := dto.SearchFilters{
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}

	if filters.SortBy == "" {
		filters.SortBy = "nickname"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "asc"
	}
	if filters.SortOrder != "asc" && filters.SortOrder != "desc" {
		filters.SortOrder = "asc"
	}

	start := time.Now()
	result, err := h.searchService.SearchPlayersWithFilters(c.Request().Context(), query, filters, page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := echo.Map{
		"query":      query,
		"filters":    filters,
		"results":    result.Results,
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"searchTime": searchTime.Milliseconds(),
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) GetPopularPlayers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := time.Now()
	result, err := h.searchService.GetPopularPlayers(c.Request().Context(), page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := echo.Map{
		"results":    result.Results,
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"searchTime": searchTime.Milliseconds(),
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) GetRecentlyActivePlayers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := time.Now()
	result, err := h.searchService.GetRecentlyActivePlayers(c.Request().Context(), page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := echo.Map{
		"results":    result.Results,
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"searchTime": searchTime.Milliseconds(),
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) SearchPlayersDebug(c echo.Context) error {
	query := c.QueryParam("query")
	searchType := c.QueryParam("searchType")

	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Query parameter is required"})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := time.Now()
	result, err := h.searchService.SearchPlayers(c.Request().Context(), query, searchType, page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return h.handleServiceError(c, err)
	}

	debugInfo := make([]echo.Map, len(result.Results))
	for i, user := range result.Results {
		debugInfo[i] = echo.Map{
			"steamID":    user.SteamID,
			"nickname":   user.Nickname,
			"avatarURL":  user.AvatarURL,
			"profileURL": user.ProfileURL,
			"createdAt":  user.CreatedAt,
			"updatedAt":  user.UpdatedAt,
			"isValid":    user.SteamID != "" && user.Nickname != "" && user.AvatarURL != "",
		}
	}

	response := echo.Map{
		"query":      query,
		"searchType": searchType,
		"searchTime": searchTime.Milliseconds(),
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"results":    result.Results,
		"debugInfo":  debugInfo,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) handleServiceError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error":   "Internal server error",
		"message": err.Error(),
	})
}
