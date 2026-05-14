package notify

import (
	"testing"
	"time"
)

func newSilencer() *Silencer {
	s := NewSilencer()
	s.now = func() time.Time { return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC) }
	return s
}

func TestSilencer_NotSilencedByDefault(t *testing.T) {
	s := newSilencer()
	if s.IsSilenced("nginx") {
		t.Fatal("expected process to not be silenced by default")
	}
}

func TestSilencer_SilencesSuppressesAlerts(t *testing.T) {
	s := newSilencer()
	future := s.now().Add(10 * time.Minute)
	s.Silence("nginx", future)
	if !s.IsSilenced("nginx") {
		t.Fatal("expected process to be silenced")
	}
}

func TestSilencer_ExpiredWindowNotSilenced(t *testing.T) {
	s := newSilencer()
	past := s.now().Add(-1 * time.Minute)
	s.Silence("nginx", past)
	if s.IsSilenced("nginx") {
		t.Fatal("expected expired silence to not suppress")
	}
}

func TestSilencer_LiftRemovesSilence(t *testing.T) {
	s := newSilencer()
	s.Silence("nginx", s.now().Add(10*time.Minute))
	s.Lift("nginx")
	if s.IsSilenced("nginx") {
		t.Fatal("expected silence to be lifted")
	}
}

func TestSilencer_LiftAllClearsAll(t *testing.T) {
	s := newSilencer()
	s.Silence("nginx", s.now().Add(10*time.Minute))
	s.Silence("redis", s.now().Add(10*time.Minute))
	s.LiftAll()
	if s.IsSilenced("nginx") || s.IsSilenced("redis") {
		t.Fatal("expected all silences to be cleared")
	}
}

func TestSilencer_IndependentPerProcess(t *testing.T) {
	s := newSilencer()
	s.Silence("nginx", s.now().Add(10*time.Minute))
	if s.IsSilenced("redis") {
		t.Fatal("silence for nginx should not affect redis")
	}
}

func TestSilencer_StatusReturnsWindow(t *testing.T) {
	s := newSilencer()
	until := s.now().Add(5 * time.Minute)
	s.Silence("nginx", until)
	got, ok := s.Status("nginx")
	if !ok {
		t.Fatal("expected status to report silenced")
	}
	if !got.Equal(until) {
		t.Fatalf("expected until %v, got %v", until, got)
	}
}
