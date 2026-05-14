package notify_test

import (
	"testing"
	"time"

	"github.com/user/procwatch/internal/notify"
)

// TestRateLimiterWithThrottle verifies that RateLimiter and Throttle
// can be composed: throttle gates per-cooldown, rate limiter gates
// the broader window burst.
func TestRateLimiterWithThrottle(t *testing.T) {
	throttle := notify.NewThrottle(20 * time.Millisecond)
	rl := notify.NewRateLimiter(time.Minute, 3)

	process := "postgres"

	send := func() bool {
		if !throttle.Allow(process) {
			return false
		}
		return rl.Allow(process)
	}

	// First alert should pass both gates.
	if !send() {
		t.Fatal("first alert should be allowed")
	}

	// Immediate second attempt blocked by throttle.
	if send() {
		t.Fatal("second immediate alert should be blocked by throttle")
	}

	// Wait for throttle cooldown.
	time.Sleep(25 * time.Millisecond)

	// Next two should pass (rate limit allows 3 total).
	if !send() {
		t.Fatal("third alert should be allowed")
	}
	time.Sleep(25 * time.Millisecond)
	if !send() {
		t.Fatal("fourth alert should be allowed")
	}

	// Rate limit exhausted (3 of 3 used).
	time.Sleep(25 * time.Millisecond)
	if send() {
		t.Fatal("should be blocked by rate limiter")
	}

	// Reset rate limiter and verify recovery.
	rl.Reset(process)
	if !send() {
		t.Fatal("should be allowed after rate limit reset")
	}
}
