package history_test

import (
	"testing"
	"time"

	"github.com/user/procwatch/internal/history"
)

func TestRecord_And_Get(t *testing.T) {
	s := history.NewStore(10)
	s.Record(history.Event{
		ProcessName: "nginx",
		Type:        history.EventCrash,
		Message:     "process exited unexpectedly",
	})

	evts := s.Get("nginx")
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Type != history.EventCrash {
		t.Errorf("expected EventCrash, got %s", evts[0].Type)
	}
	if evts[0].OccurredAt.IsZero() {
		t.Error("OccurredAt should be auto-set")
	}
}

func TestRecord_RingBuffer(t *testing.T) {
	max := 5
	s := history.NewStore(max)
	for i := 0; i < max+3; i++ {
		s.Record(history.Event{
			ProcessName: "svc",
			Type:        history.EventThreshold,
			Message:     "high cpu",
			OccurredAt:  time.Now(),
		})
	}
	if got := len(s.Get("svc")); got != max {
		t.Errorf("expected %d events (ring buffer), got %d", max, got)
	}
}

func TestAll_MultipleProcesses(t *testing.T) {
	s := history.NewStore(10)
	s.Record(history.Event{ProcessName: "nginx", Type: history.EventCrash})
	s.Record(history.Event{ProcessName: "redis", Type: history.EventThreshold})
	s.Record(history.Event{ProcessName: "redis", Type: history.EventRecovered})

	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 total events, got %d", len(all))
	}
}

func TestClear(t *testing.T) {
	s := history.NewStore(10)
	s.Record(history.Event{ProcessName: "nginx", Type: history.EventCrash})
	s.Clear("nginx")
	if got := s.Get("nginx"); len(got) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(got))
	}
}

func TestGet_UnknownProcess(t *testing.T) {
	s := history.NewStore(10)
	if got := s.Get("unknown"); got == nil || len(got) != 0 {
		t.Errorf("expected empty slice for unknown process, got %v", got)
	}
}

func TestNewStore_DefaultMax(t *testing.T) {
	s := history.NewStore(0) // should default to 50
	for i := 0; i < 55; i++ {
		s.Record(history.Event{ProcessName: "p", Type: history.EventCrash})
	}
	if got := len(s.Get("p")); got != 50 {
		t.Errorf("expected default max of 50, got %d", got)
	}
}
