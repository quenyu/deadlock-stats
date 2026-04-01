package pool

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Manager struct {
	db            *gorm.DB
	sqlDB         *sql.DB
	config        *Config
	metrics       *Metrics
	healthChecker *HealthChecker
	logger        *zap.Logger
}

func NewManager(config *Config, log *zap.Logger) (*Manager, error) {
	if log == nil {
		log, _ = zap.NewProduction()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	gormLogger := logger.Default.LogMode(logger.Silent)
	if log.Core().Enabled(zap.DebugLevel) {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		db:      db,
		sqlDB:   sqlDB,
		config:  config,
		metrics: NewMetrics(),
		logger:  log,
	}

	if err := manager.configurePool(); err != nil {
		return nil, err
	}

	manager.healthChecker = NewHealthChecker(sqlDB, config, manager.metrics, log)
	manager.healthChecker.Start()

	log.Info("database pool manager initialized",
		zap.Int("max_open_conns", config.MaxOpenConns),
		zap.Int("max_idle_conns", config.MaxIdleConns),
	)

	return manager, nil
}

// configurePool sets pool parameters
func (m *Manager) configurePool() error {
	m.sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	m.sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	m.sqlDB.SetConnMaxLifetime(m.config.ConnMaxLifetime)
	m.sqlDB.SetConnMaxIdleTime(m.config.ConnMaxIdleTime)

	m.logger.Info("connection pool configured",
		zap.Int("max_open", m.config.MaxOpenConns),
		zap.Int("max_idle", m.config.MaxIdleConns),
		zap.Duration("max_lifetime", m.config.ConnMaxLifetime),
		zap.Duration("max_idle_time", m.config.ConnMaxIdleTime),
	)

	return nil
}

// DB returns GORM instance
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// SqlDB returns sql.DB instance
func (m *Manager) SqlDB() *sql.DB {
	return m.sqlDB
}

// IsHealthy returns health status
func (m *Manager) IsHealthy() bool {
	return m.healthChecker.IsHealthy()
}

// WaitForHealthy waits for database to become healthy
func (m *Manager) WaitForHealthy(timeout time.Duration) error {
	return m.healthChecker.WaitForHealthy(timeout)
}

// GetMetrics returns current metrics
func (m *Manager) GetMetrics() Metrics {
	if m.config.EnableMetrics {
		m.updateMetrics()
	}
	return m.metrics.Snapshot()
}

// GetStats returns sql.DBStats
func (m *Manager) GetStats() sql.DBStats {
	return m.sqlDB.Stats()
}

// updateMetrics updates metrics from sql.DB stats
func (m *Manager) updateMetrics() {
	stats := m.sqlDB.Stats()
	m.metrics.Update(
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.MaxIdleClosed,
		stats.MaxLifetimeClosed,
		stats.WaitDuration,
	)
}

// LogMetrics logs current metrics
func (m *Manager) LogMetrics() {
	metrics := m.GetMetrics()

	m.logger.Info("connection pool metrics",
		zap.Int("open", metrics.OpenConnections),
		zap.Int("in_use", metrics.InUse),
		zap.Int("idle", metrics.Idle),
		zap.Int64("wait_count", metrics.WaitCount),
		zap.Duration("wait_duration", metrics.WaitDuration),
		zap.Int64("max_idle_closed", metrics.MaxIdleClosed),
		zap.Int64("max_lifetime_closed", metrics.MaxLifetimeClosed),
		zap.Int64("connection_errors", metrics.ConnectionErrors),
	)
}

// Close gracefully closes the pool
func (m *Manager) Close() error {
	m.logger.Info("closing database connection pool")

	m.healthChecker.Stop()

	if m.config.EnableMetrics {
		m.LogMetrics()
	}

	if err := m.sqlDB.Close(); err != nil {
		m.logger.Error("error closing database", zap.Error(err))
		return err
	}

	m.logger.Info("database connection pool closed")
	return nil
}
