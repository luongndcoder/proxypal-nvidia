package balancer

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter for an API key
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate int // tokens per minute
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rateLimit int) *RateLimiter {
	return &RateLimiter{
		tokens:     rateLimit,
		maxTokens:  rateLimit,
		refillRate: rateLimit,
		lastRefill: time.Now(),
	}
}

// TryAcquire attempts to acquire a token, returns true if successful
func (rl *RateLimiter) TryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// refill adds tokens based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed minutes
	tokensToAdd := int(elapsed.Minutes() * float64(rl.refillRate))

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}
}

// AvailableTokens returns the current number of available tokens
func (rl *RateLimiter) AvailableTokens() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()
	return rl.tokens
}

// TimeUntilNextToken returns the duration until the next token is available
func (rl *RateLimiter) TimeUntilNextToken() time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens > 0 {
		return 0
	}

	// Calculate time until next token
	tokensNeeded := 1
	minutesNeeded := float64(tokensNeeded) / float64(rl.refillRate)
	durationNeeded := time.Duration(minutesNeeded * float64(time.Minute))

	return durationNeeded
}
