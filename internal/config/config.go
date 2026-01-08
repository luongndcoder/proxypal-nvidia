package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	NVIDIA  NVIDIAConfig  `yaml:"nvidia"`
	Logging LoggingConfig `yaml:"logging"`
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// NVIDIAConfig contains NVIDIA API related settings
type NVIDIAConfig struct {
	BaseURL   string        `yaml:"base_url"`
	RateLimit int           `yaml:"rate_limit"`
	APIKeys   []string      `yaml:"api_keys"`
	Timeout   int           `yaml:"timeout"`
	Retry     RetryConfig   `yaml:"retry"`
}

// RetryConfig contains retry-related settings
type RetryConfig struct {
	MaxRetries   int  `yaml:"max_retries"`
	AutoFailover bool `yaml:"auto_failover"`
}

// LoggingConfig contains logging-related settings
type LoggingConfig struct {
	Level            string `yaml:"level"`
	EnableRequestLog bool   `yaml:"enable_request_log"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if len(c.NVIDIA.APIKeys) == 0 {
		return fmt.Errorf("at least one NVIDIA API key is required")
	}

	if c.NVIDIA.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be positive")
	}

	if c.NVIDIA.BaseURL == "" {
		return fmt.Errorf("NVIDIA base URL is required")
	}

	return nil
}

// GetAddress returns the server address in host:port format
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
