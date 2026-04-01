package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quenyu/deadlock-stats/internal/database/pool"
	"go.uber.org/zap"
)

type HealthHandler struct {
	poolManager *pool.Manager
	logger      *zap.Logger
}

func NewHealthHandler(poolManager *pool.Manager, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		poolManager: poolManager,
		logger:      logger,
	}
}

type HealthResponse struct {
	Status   string               `json:"status"`
	Database DatabaseHealthStatus `json:"database"`
	Metrics  *PoolMetricsResponse `json:"metrics,omitempty"`
}

type DatabaseHealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
}

type PoolMetricsResponse struct {
	OpenConnections     int    `json:"open_connections"`
	InUse               int    `json:"in_use"`
	Idle                int    `json:"idle"`
	WaitCount           int64  `json:"wait_count"`
	WaitDuration        string `json:"wait_duration"`
	MaxIdleClosed       int64  `json:"max_idle_closed"`
	MaxLifetimeClosed   int64  `json:"max_lifetime_closed"`
	ConnectionErrors    int64  `json:"connection_errors"`
	LastHealthCheck     string `json:"last_health_check"`
	HealthCheckDuration string `json:"health_check_duration"`
}

func (h *HealthHandler) HealthCheck(c echo.Context) error {
	isHealthy := h.poolManager.IsHealthy()

	response := HealthResponse{
		Status: "ok",
		Database: DatabaseHealthStatus{
			Healthy: isHealthy,
		},
	}

	if !isHealthy {
		response.Status = "degraded"
		response.Database.Message = "database connection unhealthy"
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) HealthCheckDetailed(c echo.Context) error {
	isHealthy := h.poolManager.IsHealthy()
	metrics := h.poolManager.GetMetrics()

	response := HealthResponse{
		Status: "ok",
		Database: DatabaseHealthStatus{
			Healthy: isHealthy,
		},
		Metrics: &PoolMetricsResponse{
			OpenConnections:     metrics.OpenConnections,
			InUse:               metrics.InUse,
			Idle:                metrics.Idle,
			WaitCount:           metrics.WaitCount,
			WaitDuration:        metrics.WaitDuration.String(),
			MaxIdleClosed:       metrics.MaxIdleClosed,
			MaxLifetimeClosed:   metrics.MaxLifetimeClosed,
			ConnectionErrors:    metrics.ConnectionErrors,
			LastHealthCheck:     metrics.LastHealthCheck.Format("2006-01-02T15:04:05Z07:00"),
			HealthCheckDuration: metrics.HealthCheckDuration.String(),
		},
	}

	if !isHealthy {
		response.Status = "degraded"
		response.Database.Message = "database connection unhealthy"
		return c.JSON(http.StatusServiceUnavailable, response)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) MetricsHandler(c echo.Context) error {
	stats := h.poolManager.GetStats()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	})
}
