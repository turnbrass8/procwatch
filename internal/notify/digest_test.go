package notify

import (
	"sync"
	"testing"
	"time"
)

func newDigest(window time.Duration, onFlush func(string, []DigestEntry)) *Digest {
	return NewDigest(window, onFlush)
}

func TestDigest_BatchesSingleProcess(t *testing.T) {
	var mu sync.Mutex
	flushed := map[string][]DigestEntry{}

	d := newDigest(50*time.Millisecond, func(proc string, entries []DigestEntry) {
		mu.Lock()
		flushed[proc] = entries
		mu.Unlock()
	})

	d.Add("nginx", "crash")
	d.Add("nginx", "cpu_threshold")

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(flushed["nginx"]) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(flushed["nginx"]))
	}
}

func TestDigest_PendingCount(t *testing.T) {
	d := newDigest(10*time.Second, nil)
	d.Add("redis", "crash")
	d.Add("redis", "mem_threshold")

	if got := d.PendingCount("redis"); got != 2 {
		t.Fatalf("expected 2 pending, got %d", got)
	}
}

func TestDigest_ForceFlush(t *testing.T) {
	var mu sync.Mutex
	flushed := []DigestEntry{}

	d := newDigest(10*time.Second, func(proc string, entries []DigestEntry) {
		mu.Lock()
		flushed = append(flushed, entries...)
		mu.Unlock()
	})

	d.Add("postgres", "crash")
	d.Flush("postgres")

	mu.Lock()
	defer mu.Unlock()
	if len(flushed) != 1 {
		t.Fatalf("expected 1 flushed entry, got %d", len(flushed))
	}
	if flushed[0].Reason != "crash" {
		t.Errorf("unexpected reason: %s", flushed[0].Reason)
	}
}

func TestDigest_IndependentPerProcess(t *testing.T) {
	var mu sync.Mutex
	counts := map[string]int{}

	d := newDigest(50*time.Millisecond, func(proc string, entries []DigestEntry) {
		mu.Lock()
		counts[proc] = len(entries)
		mu.Unlock()
	})

	d.Add("nginx", "crash")
	d.Add("redis", "crash")
	d.Add("redis", "cpu_threshold")

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if counts["nginx"] != 1 {
		t.Errorf("nginx: expected 1, got %d", counts["nginx"])
	}
	if counts["redis"] != 2 {
		t.Errorf("redis: expected 2, got %d", counts["redis"])
	}
}

func TestDigest_NoFlushWhenEmpty(t *testing.T) {
	called := false
	d := newDigest(10*time.Millisecond, func(_ string, _ []DigestEntry) {
		called = true
	})
	d.Flush("unknown")
	if called {
		t.Error("expected onFlush not to be called for empty process")
	}
}
