package notify_test

import (
	"testing"
	"time"

	"github.com/user/procwatch/internal/notify"
)

// TestEscalatorWithThrottle verifies that escalation and throttle work
// independently on the same process name without interfering.
func TestEscalatorWithThrottle(t *testing.T) {
	policy := notify.EscalationPolicy{
		WarningAfter:  2,
		CriticalAfter: 4,
		ResetAfter:    1 * time.Minute,
	}
	esc := notify.NewEscalator(policy)
	thr := notify.NewThrottle(30 * time.Second)

	const proc = "nginx"

	// First alert: throttle allows, escalation is Normal.
	if !thr.Allow(proc) {
		t.Fatal("first alert should be allowed by throttle")
	}
	if lvl := esc.Record(proc); lvl != notify.LevelNormal {
		t.Errorf("expected Normal on first hit, got %d", lvl)
	}

	// Throttle blocks second alert but escalation counter still advances
	// if we call Record directly (simulating a crash path that bypasses throttle).
	esc.Record(proc) // hit 2 → Warning
	esc.Record(proc) // hit 3 → Warning
	esc.Record(proc) // hit 4 → Critical

	if lvl := esc.CurrentLevel(proc); lvl != notify.LevelCritical {
		t.Errorf("expected Critical after 4 hits, got %d", lvl)
	}

	// Resetting escalation does not affect throttle.
	esc.Reset(proc)
	if esc.CurrentLevel(proc) != notify.LevelNormal {
		t.Error("escalation should be Normal after reset")
	}
	// Throttle is still blocking (cooldown not expired).
	if thr.Allow(proc) {
		t.Error("throttle should still be blocking")
	}
}
