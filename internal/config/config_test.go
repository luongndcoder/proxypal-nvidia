package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	// Create a temporary config file
	content := `
server:
  port: 8080
  host: "0.0.0.0"

nvidia:
  base_url: "https://integrate.api.nvidia.com/v1"
  rate_limit: 40
  api_keys:
    - "nvapi-key1"
    - "nvapi-key2"
  timeout: 300
  retry:
    max_retries: 3
    auto_failover: true

logging:
  level: "info"
  enable_request_log: true
`

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load config
	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Validate fields
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}

	if len(cfg.NVIDIA.APIKeys) != 2 {
		t.Errorf("Expected 2 API keys, got %d", len(cfg.NVIDIA.APIKeys))
	}

	if cfg.NVIDIA.RateLimit != 40 {
		t.Errorf("Expected rate limit 40, got %d", cfg.NVIDIA.RateLimit)
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{Port: 8080, Host: "0.0.0.0"},
				NVIDIA: NVIDIAConfig{
					BaseURL:   "https://api.nvidia.com",
					RateLimit: 40,
					APIKeys:   []string{"key1"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: Config{
				Server: ServerConfig{Port: -1, Host: "0.0.0.0"},
				NVIDIA: NVIDIAConfig{
					BaseURL:   "https://api.nvidia.com",
					RateLimit: 40,
					APIKeys:   []string{"key1"},
				},
			},
			wantErr: true,
		},
		{
			name: "no API keys",
			config: Config{
				Server: ServerConfig{Port: 8080, Host: "0.0.0.0"},
				NVIDIA: NVIDIAConfig{
					BaseURL:   "https://api.nvidia.com",
					RateLimit: 40,
					APIKeys:   []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid rate limit",
			config: Config{
				Server: ServerConfig{Port: 8080, Host: "0.0.0.0"},
				NVIDIA: NVIDIAConfig{
					BaseURL:   "https://api.nvidia.com",
					RateLimit: 0,
					APIKeys:   []string{"key1"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_GetAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	expected := "localhost:8080"
	if addr := cfg.GetAddress(); addr != expected {
		t.Errorf("GetAddress() = %s, expected %s", addr, expected)
	}
}
