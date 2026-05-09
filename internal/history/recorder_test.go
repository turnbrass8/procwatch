package history_test

import (
	"strings"
	"testing"

	"github.com/user/procwatch/internal/history"
	"github.com/user/procwatch/internal/monitor"
)

func newRecorder(t *testing.T) (*history.Recorder, *history.Store) {
	t.Helper()
	store := history.NewStore(20)
	return history.NewRecorder(store), store
}

func TestRecordCrash(t *testing.T) {
	rec, _ := newRecorder(t)
	rec.RecordCrash("nginx")
	evts := rec.RecentEvents("nginx")
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}
	if evts[0].Type != history.EventCrash {
		t.Errorf("expected EventCrash, got %s", evts[0].Type)
	}
	if !strings.Contains(evts[0].Message, "nginx") {
		t.Errorf("crash message should mention process name, got: %s", evts[0].Message)
	}
}

func TestRecordThreshold_CPUOnly(t *testing.T) {
	rec, _ := newRecorder(t)
	stats := monitor.ProcessStats{CPUPercent: 95.0, MemoryMB: 100}
	rec.RecordThreshold("redis", stats, 80.0, 512)
	evts := rec.RecentEvents("redis")
	if evts[0].Type != history.EventThreshold {
		t.Errorf("expected EventThreshold")
	}
	if !strings.Contains(evts[0].Message, "CPU") {
		t.Errorf("message should mention CPU, got: %s", evts[0].Message)
	}
}

func TestRecordThreshold_MemOnly(t *testing.T) {
	rec, _ := newRecorder(t)
	stats := monitor.ProcessStats{CPUPercent: 10.0, MemoryMB: 600}
	rec.RecordThreshold("redis", stats, 80.0, 512)
	evts := rec.RecentEvents("redis")
	if !strings.Contains(evts[0].Message, "memory") {
		t.Errorf("message should mention memory, got: %s", evts[0].Message)
	}
}

func TestRecordThreshold_Both(t *testing.T) {
	rec, _ := newRecorder(t)
	stats := monitor.ProcessStats{CPUPercent: 95.0, MemoryMB: 600}
	rec.RecordThreshold("app", stats, 80.0, 512)
	evts := rec.RecentEvents("app")
	if !strings.Contains(evts[0].Message, "CPU") || !strings.Contains(evts[0].Message, "memory") {
		t.Errorf("message should mention both CPU and memory, got: %s", evts[0].Message)
	}
}

func TestRecordRecovered(t *testing.T) {
	rec, _ := newRecorder(t)
	rec.RecordRecovered("svc")
	evts := rec.RecentEvents("svc")
	if evts[0].Type != history.EventRecovered {
		t.Errorf("expected EventRecovered, got %s", evts[0].Type)
	}
}
