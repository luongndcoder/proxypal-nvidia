package balancer

import (
	"testing"
	"time"
)

func TestRateLimiter_TryAcquire(t *testing.T) {
	// Create a rate limiter with 5 tokens
	rl := NewRateLimiter(5)

	// Should be able to acquire 5 tokens
	for i := 0; i < 5; i++ {
		if !rl.TryAcquire() {
			t.Errorf("Failed to acquire token %d", i+1)
		}
	}

	// Should fail to acquire 6th token
	if rl.TryAcquire() {
		t.Error("Should not be able to acquire token beyond limit")
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	// Create a rate limiter with 1 token per minute
	rl := NewRateLimiter(1)

	// Acquire the only token
	if !rl.TryAcquire() {
		t.Error("Failed to acquire initial token")
	}

	// Should fail to acquire another token immediately
	if rl.TryAcquire() {
		t.Error("Should not be able to acquire token immediately")
	}

	// Manually set lastRefill to 1 minute ago
	rl.lastRefill = time.Now().Add(-1 * time.Minute)

	// Should be able to acquire a token after refill
	if !rl.TryAcquire() {
		t.Error("Failed to acquire token after refill period")
	}
}

func TestRateLimiter_AvailableTokens(t *testing.T) {
	rl := NewRateLimiter(10)

	// Should have 10 tokens initially
	if tokens := rl.AvailableTokens(); tokens != 10 {
		t.Errorf("Expected 10 tokens, got %d", tokens)
	}

	// Acquire 3 tokens
	for i := 0; i < 3; i++ {
		rl.TryAcquire()
	}

	// Should have 7 tokens left
	if tokens := rl.AvailableTokens(); tokens != 7 {
		t.Errorf("Expected 7 tokens, got %d", tokens)
	}
}

func TestRateLimiter_TimeUntilNextToken(t *testing.T) {
	rl := NewRateLimiter(1)

	// Should have tokens available initially
	if duration := rl.TimeUntilNextToken(); duration != 0 {
		t.Error("Should have tokens available initially")
	}

	// Acquire the only token
	rl.TryAcquire()

	// Should need to wait for next token
	if duration := rl.TimeUntilNextToken(); duration == 0 {
		t.Error("Should need to wait for next token")
	}
}
