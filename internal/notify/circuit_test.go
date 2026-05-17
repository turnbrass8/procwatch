package notify

import (
	"testing"
	"time"
)

func newCircuit() *CircuitBreaker {
	return NewCircuitBreaker(3, 50*time.Millisecond)
}

func TestCircuit_AllowsWhenClosed(t *testing.T) {
	cb := newCircuit()
	if !cb.Allow("svc") {
		t.Fatal("expected allow when closed")
	}
}

func TestCircuit_OpensAfterMaxFailures(t *testing.T) {
	cb := newCircuit()
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}
	if cb.StateFor("svc") != StateOpen {
		t.Fatal("expected open state after max failures")
	}
	if cb.Allow("svc") {
		t.Fatal("expected block when open")
	}
}

func TestCircuit_HalfOpenAfterDuration(t *testing.T) {
	cb := newCircuit()
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow("svc") {
		t.Fatal("expected allow in half-open state")
	}
	if cb.StateFor("svc") != StateHalfOpen {
		t.Fatal("expected half-open state after duration")
	}
}

func TestCircuit_ClosesOnSuccess(t *testing.T) {
	cb := newCircuit()
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}
	time.Sleep(60 * time.Millisecond)
	cb.Allow("svc") // transitions to half-open
	cb.RecordSuccess("svc")
	if cb.StateFor("svc") != StateClosed {
		t.Fatal("expected closed after success")
	}
}

func TestCircuit_IndependentPerProcess(t *testing.T) {
	cb := newCircuit()
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svcA")
	}
	if !cb.Allow("svcB") {
		t.Fatal("svcB should not be affected by svcA failures")
	}
}

func TestCircuit_Reset_Named(t *testing.T) {
	cb := newCircuit()
	for i := 0; i < 3; i++ {
		cb.RecordFailure("svc")
	}
	cb.Reset("svc")
	if cb.StateFor("svc") != StateClosed {
		t.Fatal("expected closed after reset")
	}
}

func TestCircuit_Reset_All(t *testing.T) {
	cb := newCircuit()
	for _, name := range []string{"a", "b", "c"} {
		for i := 0; i < 3; i++ {
			cb.RecordFailure(name)
		}
	}
	cb.Reset("")
	for _, name := range []string{"a", "b", "c"} {
		if cb.StateFor(name) != StateClosed {
			t.Fatalf("expected closed for %s after full reset", name)
		}
	}
}
