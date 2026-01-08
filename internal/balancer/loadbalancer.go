package balancer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/luongndcoder/proxypal-nvidia/internal/config"
)

// APIKey represents an NVIDIA API key with its rate limiter
type APIKey struct {
	Key          string
	RateLimiter  *RateLimiter
	LastUsed     time.Time
	RequestCount atomic.Uint64
	ErrorCount   atomic.Uint64
}

// LoadBalancer manages multiple API keys and distributes requests
type LoadBalancer struct {
	apiKeys      []*APIKey
	currentIndex atomic.Uint32
	config       *config.NVIDIAConfig
	mu           sync.RWMutex
}

// NewLoadBalancer creates a new load balancer with the given API keys
func NewLoadBalancer(cfg *config.NVIDIAConfig) *LoadBalancer {
	lb := &LoadBalancer{
		apiKeys: make([]*APIKey, len(cfg.APIKeys)),
		config:  cfg,
	}

	for i, key := range cfg.APIKeys {
		lb.apiKeys[i] = &APIKey{
			Key:         key,
			RateLimiter: NewRateLimiter(cfg.RateLimit),
			LastUsed:    time.Now(),
		}
	}

	return lb
}

// GetNextKey returns the next available API key using round-robin with rate limiting
func (lb *LoadBalancer) GetNextKey() (*APIKey, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	totalKeys := len(lb.apiKeys)
	if totalKeys == 0 {
		return nil, fmt.Errorf("no API keys available")
	}

	// Try all keys starting from the current index
	startIndex := lb.currentIndex.Load()

	for i := 0; i < totalKeys; i++ {
		// Round-robin: get next index
		index := (startIndex + uint32(i)) % uint32(totalKeys)
		key := lb.apiKeys[index]

		// Try to acquire a token from this key's rate limiter
		if key.RateLimiter.TryAcquire() {
			// Update current index for next request
			lb.currentIndex.Store((index + 1) % uint32(totalKeys))

			// Update statistics
			key.LastUsed = time.Now()
			key.RequestCount.Add(1)

			return key, nil
		}
	}

	// All keys are rate limited
	return nil, fmt.Errorf("all API keys are rate limited, please wait")
}

// GetKeyWithRetry attempts to get a key with retry logic
func (lb *LoadBalancer) GetKeyWithRetry(maxRetries int) (*APIKey, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		key, err := lb.GetNextKey()
		if err == nil {
			return key, nil
		}

		lastErr = err

		// Wait a bit before retrying
		if attempt < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(attempt+1))
		}
	}

	return nil, fmt.Errorf("failed to get API key after %d retries: %w", maxRetries, lastErr)
}

// MarkKeyError increments the error count for a key
func (lb *LoadBalancer) MarkKeyError(key *APIKey) {
	if key != nil {
		key.ErrorCount.Add(1)
	}
}

// GetStats returns statistics for all API keys
func (lb *LoadBalancer) GetStats() []KeyStats {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	stats := make([]KeyStats, len(lb.apiKeys))
	for i, key := range lb.apiKeys {
		stats[i] = KeyStats{
			KeyPrefix:       MaskAPIKey(key.Key),
			RequestCount:    key.RequestCount.Load(),
			ErrorCount:      key.ErrorCount.Load(),
			AvailableTokens: key.RateLimiter.AvailableTokens(),
			LastUsed:        key.LastUsed,
		}
	}

	return stats
}

// KeyStats represents statistics for an API key
type KeyStats struct {
	KeyPrefix       string
	RequestCount    uint64
	ErrorCount      uint64
	AvailableTokens int
	LastUsed        time.Time
}

// MaskAPIKey masks an API key for security, showing only first and last few characters
func MaskAPIKey(key string) string {
	if len(key) <= 10 {
		return "***"
	}
	return key[:6] + "..." + key[len(key)-4:]
}
