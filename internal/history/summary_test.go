package history

import (
	"testing"
	"time"
)

func makeRecord(kind RecordKind, cpu, mem float64) Record {
	return Record{
		Kind:       kind,
		Timestamp:  time.Now(),
		CPUPercent: cpu,
		MemMB:      mem,
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize("svc", nil)
	if s.TotalEvents != 0 || s.CrashCount != 0 || s.ThresholdCount != 0 {
		t.Errorf("expected zero summary for empty records, got %+v", s)
	}
}

func TestSummarize_CrashesAndThresholds(t *testing.T) {
	records := []Record{
		makeRecord(KindCrash, 10, 100),
		makeRecord(KindCrash, 20, 200),
		makeRecord(KindThreshold, 90, 512),
	}
	s := Summarize("myapp", records)

	if s.ProcessName != "myapp" {
		t.Errorf("expected ProcessName=myapp, got %s", s.ProcessName)
	}
	if s.TotalEvents != 3 {
		t.Errorf("expected TotalEvents=3, got %d", s.TotalEvents)
	}
	if s.CrashCount != 2 {
		t.Errorf("expected CrashCount=2, got %d", s.CrashCount)
	}
	if s.ThresholdCount != 1 {
		t.Errorf("expected ThresholdCount=1, got %d", s.ThresholdCount)
	}
}

func TestSummarize_Averages(t *testing.T) {
	records := []Record{
		makeRecord(KindThreshold, 30, 300),
		makeRecord(KindThreshold, 60, 600),
	}
	s := Summarize("svc", records)

	if s.AvgCPU != 45.0 {
		t.Errorf("expected AvgCPU=45.0, got %f", s.AvgCPU)
	}
	if s.AvgMemMB != 450.0 {
		t.Errorf("expected AvgMemMB=450.0, got %f", s.AvgMemMB)
	}
}

func TestSummarizeAll(t *testing.T) {
	store := NewStore(10)
	store.Record("alpha", KindCrash, 5, 128)
	store.Record("beta", KindThreshold, 80, 256)
	store.Record("beta", KindCrash, 95, 512)

	summaries := SummarizeAll(store)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}

	counts := map[string]int{}
	for _, s := range summaries {
		counts[s.ProcessName] = s.TotalEvents
	}
	if counts["alpha"] != 1 {
		t.Errorf("expected alpha=1 event, got %d", counts["alpha"])
	}
	if counts["beta"] != 2 {
		t.Errorf("expected beta=2 events, got %d", counts["beta"])
	}
}
