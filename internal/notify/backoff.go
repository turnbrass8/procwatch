package notify

import (
	"math"
	"sync"
	"time"
)

// BackoffPolicy defines how delays grow between retries.
type BackoffPolicy struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultBackoffPolicy returns a sensible exponential backoff config.
func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{
		InitialDelay: 5 * time.Second,
		MaxDelay:     5 * time.Minute,
		Multiplier:   2.0,
	}
}

// BackoffTracker tracks per-process retry attempt counts and computes delays.
type BackoffTracker struct {
	mu       sync.Mutex
	policy   BackoffPolicy
	attempts map[string]int
}

// NewBackoffTracker creates a BackoffTracker with the given policy.
func NewBackoffTracker(policy BackoffPolicy) *BackoffTracker {
	return &BackoffTracker{
		policy:   policy,
		attempts: make(map[string]int),
	}
}

// Record increments the attempt counter for a process and returns the next delay.
func (b *BackoffTracker) Record(process string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.attempts[process]++
	return b.delay(b.attempts[process])
}

// Attempts returns the current attempt count for a process.
func (b *BackoffTracker) Attempts(process string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempts[process]
}

// Reset clears the attempt counter for a process.
func (b *BackoffTracker) Reset(process string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.attempts, process)
}

// ResetAll clears all attempt counters.
func (b *BackoffTracker) ResetAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.attempts = make(map[string]int)
}

func (b *BackoffTracker) delay(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}
	d := float64(b.policy.InitialDelay) * math.Pow(b.policy.Multiplier, float64(attempt-1))
	if d > float64(b.policy.MaxDelay) {
		d = float64(b.policy.MaxDelay)
	}
	return time.Duration(d)
}
