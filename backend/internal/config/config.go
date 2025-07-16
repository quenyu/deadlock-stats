package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port    string
	Timeout string
}

type DatabaseConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int
}

type RedisConfig struct {
	Host     string
	Port     string
	DB       int
	Password string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type SteamConfig struct {
	SteamAPIKey string
	Domain      string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Steam    SteamConfig
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
