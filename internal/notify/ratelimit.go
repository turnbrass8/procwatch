package notify

import (
	"sync"
	"time"
)

// RateLimiter tracks alert counts per process within a sliding window
// and prevents alert floods beyond a configurable maximum.
type RateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxAlerts int
	buckets  map[string][]time.Time
}

// NewRateLimiter creates a RateLimiter that allows at most maxAlerts
// per process within the given window duration.
func NewRateLimiter(window time.Duration, maxAlerts int) *RateLimiter {
	return &RateLimiter{
		window:    window,
		maxAlerts: maxAlerts,
		buckets:   make(map[string][]time.Time),
	}
}

// Allow returns true if an alert for the given process is permitted.
// It records the attempt and prunes stale timestamps from the window.
func (r *RateLimiter) Allow(process string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.window)

	times := r.buckets[process]
	pruned := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			pruned = append(pruned, t)
		}
	}

	if len(pruned) >= r.maxAlerts {
		r.buckets[process] = pruned
		return false
	}

	r.buckets[process] = append(pruned, now)
	return true
}

// Count returns the number of alerts recorded for a process within the window.
func (r *RateLimiter) Count(process string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-r.window)
	count := 0
	for _, t := range r.buckets[process] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears rate-limit state for a specific process or all processes
// when name is empty.
func (r *RateLimiter) Reset(process string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if process == "" {
		r.buckets = make(map[string][]time.Time)
		return
	}
	delete(r.buckets, process)
}
