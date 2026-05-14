package notify

import (
	"testing"
	"time"
)

func newDedup(window time.Duration) *Deduplicator {
	d := NewDeduplicator(window)
	return d
}

func TestDedup_FirstAlertAllowed(t *testing.T) {
	d := newDedup(5 * time.Minute)
	if d.IsDuplicate("nginx", "crash") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestDedup_SecondAlertSuppressed(t *testing.T) {
	d := newDedup(5 * time.Minute)
	d.IsDuplicate("nginx", "crash")
	if !d.IsDuplicate("nginx", "crash") {
		t.Fatal("expected duplicate alert to be suppressed")
	}
}

func TestDedup_AllowsAfterWindow(t *testing.T) {
	d := newDedup(1 * time.Second)
	fixed := time.Now()
	d.now = func() time.Time { return fixed }

	d.IsDuplicate("nginx", "crash")

	d.now = func() time.Time { return fixed.Add(2 * time.Second) }
	if d.IsDuplicate("nginx", "crash") {
		t.Fatal("expected alert after window expiry to be allowed")
	}
}

func TestDedup_DifferentReasonAllowed(t *testing.T) {
	d := newDedup(5 * time.Minute)
	d.IsDuplicate("nginx", "crash")
	if d.IsDuplicate("nginx", "high_cpu") {
		t.Fatal("expected different reason to be allowed")
	}
}

func TestDedup_IndependentPerProcess(t *testing.T) {
	d := newDedup(5 * time.Minute)
	d.IsDuplicate("nginx", "crash")
	if d.IsDuplicate("redis", "crash") {
		t.Fatal("expected different process to be allowed")
	}
}

func TestDedup_Reset_Named(t *testing.T) {
	d := newDedup(5 * time.Minute)
	d.IsDuplicate("nginx", "crash")
	d.Reset("nginx")
	if d.IsDuplicate("nginx", "crash") {
		t.Fatal("expected alert to be allowed after reset")
	}
}

func TestDedup_Reset_All(t *testing.T) {
	d := newDedup(5 * time.Minute)
	d.IsDuplicate("nginx", "crash")
	d.IsDuplicate("redis", "crash")
	d.Reset("")
	if d.IsDuplicate("nginx", "crash") {
		t.Fatal("expected nginx alert to be allowed after full reset")
	}
	if d.IsDuplicate("redis", "crash") {
		t.Fatal("expected redis alert to be allowed after full reset")
	}
}
