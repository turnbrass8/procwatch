package notify

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsUpToMax(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 3)

	for i := 0; i < 3; i++ {
		if !rl.Allow("svc") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}

	if rl.Allow("svc") {
		t.Fatal("expected deny after max alerts reached")
	}
}

func TestRateLimiter_CountReflectsWindow(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 10)
	rl.Allow("svc")
	rl.Allow("svc")

	if got := rl.Count("svc"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestRateLimiter_IndependentPerProcess(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)

	if !rl.Allow("alpha") {
		t.Fatal("alpha should be allowed")
	}
	if !rl.Allow("beta") {
		t.Fatal("beta should be allowed independently")
	}
	if rl.Allow("alpha") {
		t.Fatal("alpha should be denied after limit")
	}
}

func TestRateLimiter_Reset_Named(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)
	rl.Allow("svc")
	rl.Reset("svc")

	if !rl.Allow("svc") {
		t.Fatal("expected allow after reset")
	}
}

func TestRateLimiter_Reset_All(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)
	rl.Allow("a")
	rl.Allow("b")
	rl.Reset("")

	if !rl.Allow("a") || !rl.Allow("b") {
		t.Fatal("expected all processes reset")
	}
}

func TestRateLimiter_WindowExpiry(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 1)
	rl.Allow("svc")

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("svc") {
		t.Fatal("expected allow after window expired")
	}
}
