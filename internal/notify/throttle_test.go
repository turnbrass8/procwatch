package notify

import (
	"testing"
	"time"
)

func TestThrottle_AllowsFirstAlert(t *testing.T) {
	th := NewThrottle(5 * time.Minute)
	if !th.Allow("myapp") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestThrottle_BlocksDuringCooldown(t *testing.T) {
	th := NewThrottle(5 * time.Minute)
	th.Allow("myapp") // first call sets timestamp
	if th.Allow("myapp") {
		t.Fatal("expected second alert within cooldown to be blocked")
	}
}

func TestThrottle_AllowsAfterCooldown(t *testing.T) {
	th := NewThrottle(10 * time.Millisecond)
	th.Allow("myapp")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("myapp") {
		t.Fatal("expected alert to be allowed after cooldown expires")
	}
}

func TestThrottle_IndependentPerProcess(t *testing.T) {
	th := NewThrottle(5 * time.Minute)
	th.Allow("svcA")
	if !th.Allow("svcB") {
		t.Fatal("expected different process to be allowed independently")
	}
	if th.Allow("svcA") {
		t.Fatal("expected svcA to still be throttled")
	}
}

func TestThrottle_Reset(t *testing.T) {
	th := NewThrottle(5 * time.Minute)
	th.Allow("myapp")
	th.Reset("myapp")
	if !th.Allow("myapp") {
		t.Fatal("expected alert to be allowed after Reset")
	}
}

func TestThrottle_ResetAll(t *testing.T) {
	th := NewThrottle(5 * time.Minute)
	th.Allow("svcA")
	th.Allow("svcB")
	th.ResetAll()
	if !th.Allow("svcA") || !th.Allow("svcB") {
		t.Fatal("expected all processes to be allowed after ResetAll")
	}
}
