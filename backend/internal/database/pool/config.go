package pool

import "time"

// Config holds database connection pool configuration
type Config struct {
	// Connection settings
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string

	// Pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Health check
	HealthCheckInterval time.Duration
	EnableMetrics       bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns:        25,
		MaxIdleConns:        10,
		ConnMaxLifetime:     10 * time.Minute,
		ConnMaxIdleTime:     5 * time.Minute,
		HealthCheckInterval: 1 * time.Minute,
		EnableMetrics:       true,
	}
}

func DevelopmentConfig() *Config {
	return &Config{
		MaxOpenConns:        10,
		MaxIdleConns:        5,
		ConnMaxLifetime:     5 * time.Minute,
		ConnMaxIdleTime:     3 * time.Minute,
		HealthCheckInterval: 2 * time.Minute,
		EnableMetrics:       true,
	}
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.MaxIdleConns > c.MaxOpenConns {
		return ErrInvalidPoolConfig
	}
	if c.ConnMaxIdleTime > c.ConnMaxLifetime {
		return ErrInvalidPoolConfig
	}
	return nil
}

// DSN builds database connection string
func (c *Config) DSN() string {
	return buildDSN(c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func buildDSN(host, port, user, password, dbname, sslmode string) string {
	return "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" sslmode=" + sslmode
}
