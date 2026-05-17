package notify

import (
	"testing"
	"time"
)

func newBackoff() *BackoffTracker {
	return NewBackoffTracker(BackoffPolicy{
		InitialDelay: 1 * time.Second,
		MaxDelay:     16 * time.Second,
		Multiplier:   2.0,
	})
}

func TestBackoff_FirstAttemptReturnsInitialDelay(t *testing.T) {
	b := newBackoff()
	d := b.Record("svc")
	if d != 1*time.Second {
		t.Errorf("expected 1s, got %v", d)
	}
}

func TestBackoff_DelayGrowsExponentially(t *testing.T) {
	b := newBackoff()
	expected := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second}
	for i, want := range expected {
		got := b.Record("svc")
		if got != want {
			t.Errorf("attempt %d: expected %v, got %v", i+1, want, got)
		}
	}
}

func TestBackoff_CapsAtMaxDelay(t *testing.T) {
	b := newBackoff()
	var last time.Duration
	for i := 0; i < 10; i++ {
		last = b.Record("svc")
	}
	if last != 16*time.Second {
		t.Errorf("expected max delay 16s, got %v", last)
	}
}

func TestBackoff_ResetClearsAttempts(t *testing.T) {
	b := newBackoff()
	b.Record("svc")
	b.Record("svc")
	b.Reset("svc")
	if b.Attempts("svc") != 0 {
		t.Errorf("expected 0 attempts after reset, got %d", b.Attempts("svc"))
	}
	d := b.Record("svc")
	if d != 1*time.Second {
		t.Errorf("expected initial delay after reset, got %v", d)
	}
}

func TestBackoff_IndependentPerProcess(t *testing.T) {
	b := newBackoff()
	b.Record("svc-a")
	b.Record("svc-a")
	d := b.Record("svc-b")
	if d != 1*time.Second {
		t.Errorf("svc-b should start fresh, got %v", d)
	}
}

func TestBackoff_ResetAll(t *testing.T) {
	b := newBackoff()
	b.Record("svc-a")
	b.Record("svc-b")
	b.ResetAll()
	if b.Attempts("svc-a") != 0 || b.Attempts("svc-b") != 0 {
		t.Error("expected all attempts cleared after ResetAll")
	}
}

func TestDefaultBackoffPolicy(t *testing.T) {
	p := DefaultBackoffPolicy()
	if p.InitialDelay <= 0 {
		t.Error("InitialDelay should be positive")
	}
	if p.MaxDelay < p.InitialDelay {
		t.Error("MaxDelay should be >= InitialDelay")
	}
	if p.Multiplier <= 1.0 {
		t.Error("Multiplier should be > 1")
	}
}
