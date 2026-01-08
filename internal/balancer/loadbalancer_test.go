package balancer

import (
	"testing"

	"github.com/luongndcoder/proxypal-nvidia/internal/config"
)

func TestLoadBalancer_NewLoadBalancer(t *testing.T) {
	cfg := &config.NVIDIAConfig{
		APIKeys:   []string{"key1", "key2", "key3"},
		RateLimit: 40,
	}

	lb := NewLoadBalancer(cfg)

	if len(lb.apiKeys) != 3 {
		t.Errorf("Expected 3 API keys, got %d", len(lb.apiKeys))
	}
}

func TestLoadBalancer_GetNextKey(t *testing.T) {
	cfg := &config.NVIDIAConfig{
		APIKeys:   []string{"key1", "key2", "key3"},
		RateLimit: 40,
	}

	lb := NewLoadBalancer(cfg)

	// Should get keys in round-robin fashion
	keys := make([]string, 6)
	for i := 0; i < 6; i++ {
		key, err := lb.GetNextKey()
		if err != nil {
			t.Fatalf("Failed to get key: %v", err)
		}
		keys[i] = key.Key
	}

	// Check round-robin pattern: key1, key2, key3, key1, key2, key3
	expected := []string{"key1", "key2", "key3", "key1", "key2", "key3"}
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Key %d: expected %s, got %s", i, expected[i], key)
		}
	}
}

func TestLoadBalancer_GetStats(t *testing.T) {
	cfg := &config.NVIDIAConfig{
		APIKeys:   []string{"key1", "key2"},
		RateLimit: 40,
	}

	lb := NewLoadBalancer(cfg)

	// Get a key to increment request count
	key, err := lb.GetNextKey()
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if key.RequestCount.Load() != 1 {
		t.Errorf("Expected request count 1, got %d", key.RequestCount.Load())
	}

	// Get stats
	stats := lb.GetStats()
	if len(stats) != 2 {
		t.Errorf("Expected 2 stats entries, got %d", len(stats))
	}

	// First key should have 1 request
	if stats[0].RequestCount != 1 {
		t.Errorf("Expected request count 1 for first key, got %d", stats[0].RequestCount)
	}

	// Second key should have 0 requests
	if stats[1].RequestCount != 0 {
		t.Errorf("Expected request count 0 for second key, got %d", stats[1].RequestCount)
	}
}

func TestLoadBalancer_MarkKeyError(t *testing.T) {
	cfg := &config.NVIDIAConfig{
		APIKeys:   []string{"key1"},
		RateLimit: 40,
	}

	lb := NewLoadBalancer(cfg)

	key := lb.apiKeys[0]
	initialErrors := key.ErrorCount.Load()

	lb.MarkKeyError(key)

	if key.ErrorCount.Load() != initialErrors+1 {
		t.Errorf("Expected error count to increase by 1")
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"nvapi-1234567890abcdef", "nvapi-...cdef"},
		{"short", "***"},
		{"nvapi-very-long-api-key-here", "nvapi-...here"},
	}

	for _, test := range tests {
		result := MaskAPIKey(test.input)
		if result != test.expected {
			t.Errorf("MaskAPIKey(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}
