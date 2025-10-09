package ratelimit

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Manager struct {
	config     *Config
	middleware *Middleware
	limiter    Limiter
	logger     *zap.Logger
}

type ManagerConfig struct {
	Config      *Config
	RedisClient *redis.Client
	Logger      *zap.Logger
}

func NewManager(cfg *ManagerConfig) (*Manager, error) {
	if cfg.Config == nil {
		cfg.Config = DefaultConfig()
	}

	if cfg.Logger == nil {
		cfg.Logger, _ = zap.NewProduction()
	}

	// Validate configuration
	if err := cfg.Config.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	// Create limiter based on configuration
	var limiter Limiter
	if cfg.Config.UseRedis && cfg.RedisClient != nil {
		ttl := cfg.Config.RedisKeyTTL
		if ttl == 0 {
			ttl = time.Minute
		}
		limiter = NewRedisLimiter(cfg.RedisClient, ttl)
		cfg.Logger.Info("using Redis-based rate limiter")
	} else {
		ttl := cfg.Config.RedisKeyTTL
		if ttl == 0 {
			ttl = 5 * time.Minute
		}
		limiter = NewMemoryLimiter(cfg.Config.Burst, ttl)
		cfg.Logger.Info("using in-memory rate limiter")
	}

	// Create middleware
	middleware, err := NewMiddleware(cfg.Config, limiter, cfg.Logger)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info("rate limiting manager initialized",
		zap.String("strategy", string(cfg.Config.Strategy)),
		zap.Int("requests_per_second", cfg.Config.RequestsPerSecond),
		zap.Int("burst", cfg.Config.Burst),
		zap.Bool("use_redis", cfg.Config.UseRedis),
	)

	return &Manager{
		config:     cfg.Config,
		middleware: middleware,
		limiter:    limiter,
		logger:     cfg.Logger,
	}, nil
}

// Middleware returns the Echo middleware function
func (m *Manager) Middleware() echo.MiddlewareFunc {
	return m.middleware.Handler()
}

// SetOnLimitReached sets callback for when rate limit is exceeded
func (m *Manager) SetOnLimitReached(callback LimitReachedCallback) {
	m.middleware.SetOnLimitReached(callback)
}

// Limiter returns the underlying limiter
func (m *Manager) Limiter() Limiter {
	return m.limiter
}

// Config returns the configuration
func (m *Manager) Config() *Config {
	return m.config
}

// UpdateConfig updates runtime configuration (be careful with this)
func (m *Manager) UpdateConfig(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	m.config = cfg
	m.middleware.config = cfg
	m.logger.Info("rate limiter configuration updated")
	return nil
}

// Close closes the manager and underlying resources
func (m *Manager) Close() error {
	m.logger.Info("closing rate limiting manager")

	if err := m.middleware.Close(); err != nil {
		m.logger.Error("error closing middleware", zap.Error(err))
		return err
	}

	if err := m.limiter.Close(); err != nil {
		m.logger.Error("error closing limiter", zap.Error(err))
		return err
	}

	m.logger.Info("rate limiting manager closed")
	return nil
}
