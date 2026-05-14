package notify

import (
	"sync"
	"time"
)

// Throttle prevents duplicate alerts from firing too frequently
// for the same process within a configurable cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSent map[string]time.Time
}

// NewThrottle creates a Throttle with the given cooldown duration.
func NewThrottle(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}
}

// Allow returns true if an alert for the given process name is permitted,
// i.e. no alert has been sent within the cooldown window. It records the
// current time as the last-sent timestamp when returning true.
func (t *Throttle) Allow(process string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastSent[process]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSent[process] = now
	return true
}

// Reset clears the last-sent record for a specific process.
func (t *Throttle) Reset(process string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, process)
}

// ResetAll clears all throttle state.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSent = make(map[string]time.Time)
}
