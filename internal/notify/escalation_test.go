package notify

import (
	"testing"
	"time"
)

func defaultPolicy() EscalationPolicy {
	return EscalationPolicy{
		WarningAfter:  2,
		CriticalAfter: 4,
		ResetAfter:    10 * time.Minute,
	}
}

func TestEscalator_NormalOnFirstHit(t *testing.T) {
	e := NewEscalator(defaultPolicy())
	level := e.Record("svc")
	if level != LevelNormal {
		t.Errorf("expected Normal, got %d", level)
	}
}

func TestEscalator_WarningAfterThreshold(t *testing.T) {
	e := NewEscalator(defaultPolicy())
	e.Record("svc")
	level := e.Record("svc")
	if level != LevelWarning {
		t.Errorf("expected Warning, got %d", level)
	}
}

func TestEscalator_CriticalAfterThreshold(t *testing.T) {
	e := NewEscalator(defaultPolicy())
	for i := 0; i < 4; i++ {
		e.Record("svc")
	}
	level := e.CurrentLevel("svc")
	if level != LevelCritical {
		t.Errorf("expected Critical, got %d", level)
	}
}

func TestEscalator_ResetClearsState(t *testing.T) {
	e := NewEscalator(defaultPolicy())
	for i := 0; i < 4; i++ {
		e.Record("svc")
	}
	e.Reset("svc")
	if e.CurrentLevel("svc") != LevelNormal {
		t.Error("expected Normal after reset")
	}
}

func TestEscalator_IndependentPerProcess(t *testing.T) {
	e := NewEscalator(defaultPolicy())
	for i := 0; i < 4; i++ {
		e.Record("svc-a")
	}
	level := e.Record("svc-b")
	if level != LevelNormal {
		t.Errorf("svc-b should be Normal, got %d", level)
	}
}

func TestEscalator_ResetsAfterQuietPeriod(t *testing.T) {
	policy := EscalationPolicy{
		WarningAfter:  2,
		CriticalAfter: 4,
		ResetAfter:    50 * time.Millisecond,
	}
	e := NewEscalator(policy)
	for i := 0; i < 4; i++ {
		e.Record("svc")
	}
	time.Sleep(60 * time.Millisecond)
	level := e.Record("svc")
	if level != LevelNormal {
		t.Errorf("expected Normal after quiet period, got %d", level)
	}
}

func TestDefaultEscalationPolicy(t *testing.T) {
	p := DefaultEscalationPolicy()
	if p.WarningAfter <= 0 || p.CriticalAfter <= 0 || p.ResetAfter <= 0 {
		t.Error("default policy fields must be positive")
	}
}
