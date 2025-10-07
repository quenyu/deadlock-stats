package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/dto"
	cErrors "github.com/quenyu/deadlock-stats/internal/errors"
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

	if query == "" {
		return ErrorHandler(cErrors.ErrInvalidQuery, c)
	}

	page, pageSize := parsePaginationParams(c, 1, 10)

	result, err := h.searchService.SearchPlayers(c.Request().Context(), query, searchType, page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		h.logger.Error("SearchPlayers error", zap.Error(err))
		return ErrorHandler(err, c)
	}

	if result.TotalCount == 0 {
		return ErrorHandler(cErrors.ErrNoSearchResults, c)
	}

	response := echo.Map{
		"results":    result.Results,
		"totalCount": result.TotalCount,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"totalPages": result.TotalPages,
		"searchType": searchType,
		"searchTime": searchTime.Milliseconds(),
		"query":      query,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) SearchPlayersAutocomplete(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return ErrorHandler(cErrors.ErrInvalidQuery, c)
	}

	limit := parseLimit(c, 10, 50)
	start := time.Now()

	users, err := h.searchService.SearchPlayersWithAutocomplete(c.Request().Context(), query, limit)
	searchTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if len(users) == 0 {
		return ErrorHandler(cErrors.ErrNoSearchResults, c)
	}

	response := echo.Map{
		"query":      query,
		"limit":      limit,
		"searchTime": searchTime.Milliseconds(),
		"totalFound": len(users),
		"results":    users,
	}

	h.logger.Info("Autocomplete search completed",
		zap.String("query", query),
		zap.Int("totalFound", len(users)),
	)

	return c.JSON(http.StatusOK, response)
}

func (h *PlayerSearchHandler) SearchPlayersWithFilters(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return ErrorHandler(cErrors.ErrInvalidQuery, c)
	}

	page, pageSize := parsePaginationParams(c, 1, 20)
	filters := parseSearchFilters(c)

	start := time.Now()
	result, err := h.searchService.SearchPlayersWithFilters(c.Request().Context(), query, filters, page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if result.TotalCount == 0 {
		return ErrorHandler(cErrors.ErrNoSearchResults, c)
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
	page, pageSize := parsePaginationParams(c, 1, 10)

	start := time.Now()
	result, err := h.searchService.GetPopularPlayers(c.Request().Context(), page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
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
	page, pageSize := parsePaginationParams(c, 1, 10)
	start := time.Now()

	result, err := h.searchService.GetRecentlyActivePlayers(c.Request().Context(), page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
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
		return ErrorHandler(cErrors.ErrInvalidQuery, c)
	}

	page, pageSize := parsePaginationParams(c, 1, 10)

	start := time.Now()
	result, err := h.searchService.SearchPlayers(c.Request().Context(), query, searchType, page, pageSize)
	searchTime := time.Since(start)

	if err != nil {
		return ErrorHandler(err, c)
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

func parsePaginationParams(c echo.Context, defaultPage, defaultSize int) (int, int) {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultSize
	}
	return page, pageSize
}

func parseLimit(c echo.Context, defaultLimit, maxLimit int) int {
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		return defaultLimit
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= maxLimit {
		return l
	}
	return defaultLimit
}

func parseSearchFilters(c echo.Context) dto.SearchFilters {
	filters := dto.SearchFilters{
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}
	if filters.SortBy == "" {
		filters.SortBy = "nickname"
	}
	if filters.SortOrder != "asc" && filters.SortOrder != "desc" {
		filters.SortOrder = "asc"
	}
	return filters
}
