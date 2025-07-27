package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Steam    SteamConfig    `mapstructure:"steam"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	App      AppConfig      `mapstructure:"app"`
	API      APIConfig      `mapstructure:"api"`
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
	Version string `mapstructure:"version"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	SSLMode        string `mapstructure:"sslmode"`
	MaxConnections int    `mapstructure:"max_connections"`
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
