package pool

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"go.uber.org/zap"
)

type HealthChecker struct {
	db        *sql.DB
	config    *Config
	metrics   *Metrics
	logger    *zap.Logger
	stopChan  chan struct{}
	mu        sync.RWMutex
	isHealthy bool
}

func NewHealthChecker(db *sql.DB, config *Config, metrics *Metrics, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		db:        db,
		config:    config,
		metrics:   metrics,
		logger:    logger,
		stopChan:  make(chan struct{}),
		isHealthy: false,
	}
}

// Start begins periodic health checks
func (h *HealthChecker) Start() {
	if h.config.HealthCheckInterval == 0 {
		h.logger.Info("health checks disabled")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = h.Check(ctx)
	cancel()

	go h.run()
	h.logger.Info("health checker started", zap.Duration("interval", h.config.HealthCheckInterval))
}

// Stop stops health checks
func (h *HealthChecker) Stop() {
	close(h.stopChan)
	h.logger.Info("health checker stopped")
}

// IsHealthy returns current health status
func (h *HealthChecker) IsHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isHealthy
}

// Check performs a health check
func (h *HealthChecker) Check(ctx context.Context) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		h.metrics.UpdateHealthCheck(duration)
	}()

	if err := h.db.PingContext(ctx); err != nil {
		h.setHealthy(false)
		h.metrics.IncrementErrors()
		h.logger.Error("health check failed", zap.Error(err))
		return err
	}

	h.setHealthy(true)
	return nil
}

// WaitForHealthy waits for database to become healthy
func (h *HealthChecker) WaitForHealthy(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		case <-ticker.C:
			if err := h.Check(context.Background()); err == nil {
				return nil
			}
		}
	}
}

// run executes periodic health checks
func (h *HealthChecker) run() {
	ticker := time.NewTicker(h.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := h.Check(ctx); err != nil {
				h.logger.Warn("periodic health check failed", zap.Error(err))
			} else {
				h.logger.Debug("health check passed",
					zap.Int("open", h.metrics.OpenConnections),
					zap.Int("idle", h.metrics.Idle),
				)
			}
			cancel()

		case <-h.stopChan:
			return
		}
	}
}

// setHealthy updates health status
func (h *HealthChecker) setHealthy(healthy bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.isHealthy = healthy
}
