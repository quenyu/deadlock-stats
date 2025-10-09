package ratelimit

import (
	"fmt"
	"time"
)

type Strategy string

const (
	StrategyIP        Strategy = "ip"
	StrategyUser      Strategy = "user"
	StrategyIPAndUser Strategy = "ip_and_user"
	StrategyEndpoint  Strategy = "endpoint"
	StrategyCustom    Strategy = "custom"
)

type Config struct {
	// General settings
	Enabled           bool
	Strategy          Strategy
	RequestsPerSecond int
	Burst             int

	// Redis settings
	UseRedis    bool
	RedisKeyTTL time.Duration

	// Per-endpoint limits (endpoint -> requests per second)
	PerEndpoint map[string]int

	// IP whitelist (no rate limiting)
	Whitelist []string

	// Trusted proxies for real IP detection
	TrustedProxies []string

	// Custom key function
	KeyFunc KeyExtractor
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		Strategy:          StrategyIP,
		RequestsPerSecond: 100,
		Burst:             200,
		UseRedis:          true,
		RedisKeyTTL:       time.Minute,
		PerEndpoint:       make(map[string]int),
		Whitelist:         []string{"127.0.0.1", "::1"},
		TrustedProxies:    []string{"127.0.0.1"},
	}
}

func DevelopmentConfig() *Config {
	return &Config{
		Enabled:           true,
		Strategy:          StrategyIP,
		RequestsPerSecond: 1000,
		Burst:             2000,
		UseRedis:          false,
		PerEndpoint:       make(map[string]int),
		Whitelist:         []string{"127.0.0.1", "::1"},
		TrustedProxies:    []string{"127.0.0.1"},
	}
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.RequestsPerSecond <= 0 {
		return fmt.Errorf("requests_per_second must be positive")
	}
	if c.Burst <= 0 {
		return fmt.Errorf("burst must be positive")
	}
	if c.Burst < c.RequestsPerSecond {
		return fmt.Errorf("burst must be >= requests_per_second")
	}
	return nil
}

// GetLimit returns the limit for a specific endpoint
func (c *Config) GetLimit(endpoint string) int {
	if limit, ok := c.PerEndpoint[endpoint]; ok {
		return limit
	}
	return c.RequestsPerSecond
}

// IsWhitelisted checks if IP is whitelisted
func (c *Config) IsWhitelisted(ip string) bool {
	for _, whitelistedIP := range c.Whitelist {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}
