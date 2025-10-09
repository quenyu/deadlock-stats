package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Steam     SteamConfig     `mapstructure:"steam"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	App       AppConfig       `mapstructure:"app"`
	API       APIConfig       `mapstructure:"api"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Security  SecurityConfig  `mapstructure:"security"`
}

type APIConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	MaxRetries      int           `mapstructure:"max_retries"`
	ConnectionPool  int           `mapstructure:"connection_pool"`
	IdleConnTimeout time.Duration `mapstructure:"idle_conn_timeout"`
	CacheTTL        time.Duration `mapstructure:"cache_ttl"`
	PartialCacheTTL time.Duration `mapstructure:"partial_cache_ttl"`
	EnableMetrics   bool          `mapstructure:"enable_metrics"`
	EnableRetry     bool          `mapstructure:"enable_retry"`
}

type AppConfig struct {
	Version   string `mapstructure:"version"`
	ClientURL string `mapstructure:"client_url"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host           string     `mapstructure:"host"`
	Port           string     `mapstructure:"port"`
	User           string     `mapstructure:"user"`
	Password       string     `mapstructure:"password"`
	Name           string     `mapstructure:"name"`
	SSLMode        string     `mapstructure:"sslmode"`
	MaxConnections int        `mapstructure:"max_connections"`
	Pool           PoolConfig `mapstructure:"pool"`
}

// PoolConfig holds database connection pool configuration
type PoolConfig struct {
	MaxOpenConns        int           `mapstructure:"max_open_conns"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime     time.Duration `mapstructure:"conn_max_idle_time"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`

	// EnableMetrics enables connection pool metrics collection
	EnableMetrics bool `mapstructure:"enable_metrics"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type SteamConfig struct {
	RedirectURL string `mapstructure:"domain"`
	APIKey      string `mapstructure:"steam_api_key"`
}

type RateLimitConfig struct {
	Enabled           bool           `mapstructure:"enabled"`
	Strategy          string         `mapstructure:"strategy"`
	RequestsPerSecond int            `mapstructure:"requests_per_second"`
	Burst             int            `mapstructure:"burst"`
	UseRedis          bool           `mapstructure:"use_redis"`
	RedisKeyTTL       time.Duration  `mapstructure:"redis_key_ttl"`
	PerEndpoint       map[string]int `mapstructure:"per_endpoint"`
	Whitelist         []string       `mapstructure:"whitelist"`
	TrustedProxies    []string       `mapstructure:"trusted_proxies"`
}

type SecurityConfig struct {
	// HSTS
	HSTSEnabled           bool `mapstructure:"hsts_enabled"`
	HSTSMaxAge            int  `mapstructure:"hsts_max_age"`
	HSTSIncludeSubdomains bool `mapstructure:"hsts_include_subdomains"`
	HSTSPreload           bool `mapstructure:"hsts_preload"`

	// CSP
	CSPEnabled    bool `mapstructure:"csp_enabled"`
	CSPReportOnly bool `mapstructure:"csp_report_only"`

	// XSS
	XSSProtection      bool   `mapstructure:"xss_protection"`
	XFrameOptions      string `mapstructure:"x_frame_options"`
	ContentTypeNoSniff bool   `mapstructure:"content_type_no_sniff"`

	// Other
	ReferrerPolicy    string `mapstructure:"referrer_policy"`
	PermissionsPolicy string `mapstructure:"permissions_policy"`

	// CSRF
	CSRFEnabled        bool   `mapstructure:"csrf_enabled"`
	CSRFCookieSecure   bool   `mapstructure:"csrf_cookie_secure"`
	CSRFCookieSameSite string `mapstructure:"csrf_cookie_same_site"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
