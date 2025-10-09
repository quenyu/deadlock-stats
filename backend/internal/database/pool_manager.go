package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/quenyu/deadlock-stats/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PoolManager manages database connection pool with health checks and metrics
type PoolManager struct {
	db              *gorm.DB
	sqlDB           *sql.DB
	config          *config.DatabaseConfig
	logger          *zap.Logger
	metrics         *PoolMetrics
	healthCheckStop chan struct{}
	mu              sync.RWMutex
	isHealthy       bool
}

// PoolMetrics holds connection pool metrics
type PoolMetrics struct {
	mu                  sync.RWMutex
	OpenConnections     int
	InUse               int
	Idle                int
	WaitCount           int64
	WaitDuration        time.Duration
	MaxIdleClosed       int64
	MaxLifetimeClosed   int64
	ConnectionErrors    int64
	LastHealthCheck     time.Time
	HealthCheckDuration time.Duration
}

func NewPoolManager(cfg *config.DatabaseConfig, log *zap.Logger) (*PoolManager, error) {
	if log == nil {
		log, _ = zap.NewProduction()
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	gormLogger := logger.Default.LogMode(logger.Silent)
	if log.Core().Enabled(zap.DebugLevel) {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	pm := &PoolManager{
		db:              db,
		sqlDB:           sqlDB,
		config:          cfg,
		logger:          log,
		metrics:         &PoolMetrics{},
		healthCheckStop: make(chan struct{}),
		isHealthy:       false,
	}

	if err := pm.configurePool(); err != nil {
		return nil, fmt.Errorf("failed to configure pool: %w", err)
	}

	if err := pm.checkHealth(context.Background()); err != nil {
		log.Warn("initial health check failed", zap.Error(err))
	}

	if cfg.Pool.HealthCheckInterval > 0 {
		go pm.runHealthChecks()
	}

	log.Info("database pool manager initialized",
		zap.Int("max_open_conns", cfg.Pool.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Pool.MaxIdleConns),
		zap.Duration("conn_max_lifetime", cfg.Pool.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", cfg.Pool.ConnMaxIdleTime),
	)

	return pm, nil
}

func (pm *PoolManager) configurePool() error {
	poolCfg := pm.config.Pool

	if poolCfg.MaxOpenConns == 0 {
		poolCfg.MaxOpenConns = 25
	}
	if poolCfg.MaxIdleConns == 0 {
		poolCfg.MaxIdleConns = 10
	}
	if poolCfg.ConnMaxLifetime == 0 {
		poolCfg.ConnMaxLifetime = 10 * time.Minute
	}
	if poolCfg.ConnMaxIdleTime == 0 {
		poolCfg.ConnMaxIdleTime = 5 * time.Minute
	}

	if poolCfg.MaxIdleConns > poolCfg.MaxOpenConns {
		return fmt.Errorf("max_idle_conns (%d) cannot exceed max_open_conns (%d)",
			poolCfg.MaxIdleConns, poolCfg.MaxOpenConns)
	}

	if poolCfg.ConnMaxIdleTime > poolCfg.ConnMaxLifetime {
		pm.logger.Warn("conn_max_idle_time is greater than conn_max_lifetime, adjusting",
			zap.Duration("idle_time", poolCfg.ConnMaxIdleTime),
			zap.Duration("lifetime", poolCfg.ConnMaxLifetime),
		)
		poolCfg.ConnMaxIdleTime = poolCfg.ConnMaxLifetime - time.Minute
	}

	pm.sqlDB.SetMaxOpenConns(poolCfg.MaxOpenConns)
	pm.sqlDB.SetMaxIdleConns(poolCfg.MaxIdleConns)
	pm.sqlDB.SetConnMaxLifetime(poolCfg.ConnMaxLifetime)
	pm.sqlDB.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)

	pm.logger.Info("connection pool configured",
		zap.Int("max_open_conns", poolCfg.MaxOpenConns),
		zap.Int("max_idle_conns", poolCfg.MaxIdleConns),
		zap.Duration("conn_max_lifetime", poolCfg.ConnMaxLifetime),
		zap.Duration("conn_max_idle_time", poolCfg.ConnMaxIdleTime),
	)

	return nil
}

func (pm *PoolManager) DB() *gorm.DB {
	return pm.db
}

func (pm *PoolManager) SqlDB() *sql.DB {
	return pm.sqlDB
}

func (pm *PoolManager) IsHealthy() bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.isHealthy
}

func (pm *PoolManager) checkHealth(ctx context.Context) error {
	start := time.Now()
	defer func() {
		pm.metrics.mu.Lock()
		pm.metrics.LastHealthCheck = time.Now()
		pm.metrics.HealthCheckDuration = time.Since(start)
		pm.metrics.mu.Unlock()
	}()

	if err := pm.sqlDB.PingContext(ctx); err != nil {
		pm.mu.Lock()
		pm.isHealthy = false
		pm.mu.Unlock()

		pm.metrics.mu.Lock()
		pm.metrics.ConnectionErrors++
		pm.metrics.mu.Unlock()

		pm.logger.Error("database health check failed", zap.Error(err))
		return err
	}

	pm.mu.Lock()
	pm.isHealthy = true
	pm.mu.Unlock()

	if pm.config.Pool.EnableMetrics {
		pm.updateMetrics()
	}

	return nil
}

func (pm *PoolManager) updateMetrics() {
	stats := pm.sqlDB.Stats()

	pm.metrics.mu.Lock()
	defer pm.metrics.mu.Unlock()

	pm.metrics.OpenConnections = stats.OpenConnections
	pm.metrics.InUse = stats.InUse
	pm.metrics.Idle = stats.Idle
	pm.metrics.WaitCount = stats.WaitCount
	pm.metrics.WaitDuration = stats.WaitDuration
	pm.metrics.MaxIdleClosed = stats.MaxIdleClosed
	pm.metrics.MaxLifetimeClosed = stats.MaxLifetimeClosed
}

func (pm *PoolManager) GetMetrics() PoolMetrics {
	pm.metrics.mu.RLock()
	defer pm.metrics.mu.RUnlock()

	return PoolMetrics{
		OpenConnections:     pm.metrics.OpenConnections,
		InUse:               pm.metrics.InUse,
		Idle:                pm.metrics.Idle,
		WaitCount:           pm.metrics.WaitCount,
		WaitDuration:        pm.metrics.WaitDuration,
		MaxIdleClosed:       pm.metrics.MaxIdleClosed,
		MaxLifetimeClosed:   pm.metrics.MaxLifetimeClosed,
		ConnectionErrors:    pm.metrics.ConnectionErrors,
		LastHealthCheck:     pm.metrics.LastHealthCheck,
		HealthCheckDuration: pm.metrics.HealthCheckDuration,
	}
}

func (pm *PoolManager) runHealthChecks() {
	interval := pm.config.Pool.HealthCheckInterval
	if interval == 0 {
		interval = 1 * time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	pm.logger.Info("starting health check routine", zap.Duration("interval", interval))

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := pm.checkHealth(ctx); err != nil {
				pm.logger.Warn("periodic health check failed",
					zap.Error(err),
					zap.Int("open_connections", pm.metrics.OpenConnections),
				)
			} else {
				pm.logger.Debug("health check passed",
					zap.Int("open_connections", pm.metrics.OpenConnections),
					zap.Int("in_use", pm.metrics.InUse),
					zap.Int("idle", pm.metrics.Idle),
				)
			}
			cancel()

		case <-pm.healthCheckStop:
			pm.logger.Info("stopping health check routine")
			return
		}
	}
}

func (pm *PoolManager) LogMetrics() {
	metrics := pm.GetMetrics()

	pm.logger.Info("connection pool metrics",
		zap.Int("open_connections", metrics.OpenConnections),
		zap.Int("in_use", metrics.InUse),
		zap.Int("idle", metrics.Idle),
		zap.Int64("wait_count", metrics.WaitCount),
		zap.Duration("wait_duration", metrics.WaitDuration),
		zap.Int64("max_idle_closed", metrics.MaxIdleClosed),
		zap.Int64("max_lifetime_closed", metrics.MaxLifetimeClosed),
		zap.Int64("connection_errors", metrics.ConnectionErrors),
		zap.Time("last_health_check", metrics.LastHealthCheck),
		zap.Duration("health_check_duration", metrics.HealthCheckDuration),
	)
}

func (pm *PoolManager) Close() error {
	pm.logger.Info("closing database connection pool")

	close(pm.healthCheckStop)

	if pm.config.Pool.EnableMetrics {
		pm.LogMetrics()
	}

	if err := pm.sqlDB.Close(); err != nil {
		pm.logger.Error("error closing database", zap.Error(err))
		return err
	}

	pm.logger.Info("database connection pool closed successfully")
	return nil
}

func (pm *PoolManager) WaitForHealthy(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database to become healthy")
		case <-ticker.C:
			if err := pm.checkHealth(context.Background()); err == nil {
				return nil
			}
		}
	}
}

func (pm *PoolManager) GetStats() sql.DBStats {
	return pm.sqlDB.Stats()
}
