package notify

import (
	"sync"
	"time"
)

// DigestEntry holds a single alert event for batching.
type DigestEntry struct {
	Process string
	Reason  string
	At      time.Time
}

// Digest batches alert events per process and flushes them after a window.
type Digest struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string][]DigestEntry
	timers  map[string]*time.Timer
	onFlush func(process string, entries []DigestEntry)
}

// NewDigest creates a Digest that calls onFlush when the window expires for a process.
func NewDigest(window time.Duration, onFlush func(process string, entries []DigestEntry)) *Digest {
	return &Digest{
		window:  window,
		entries: make(map[string][]DigestEntry),
		timers:  make(map[string]*time.Timer),
		onFlush: onFlush,
	}
}

// Add records an alert entry for the given process. If no flush timer is
// running for the process, one is started.
func (d *Digest) Add(process, reason string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.entries[process] = append(d.entries[process], DigestEntry{
		Process: process,
		Reason:  reason,
		At:      time.Now(),
	})

	if _, ok := d.timers[process]; !ok {
		d.timers[process] = time.AfterFunc(d.window, func() {
			d.flush(process)
		})
	}
}

// Flush forces an immediate flush for the given process.
func (d *Digest) Flush(process string) {
	d.mu.Lock()
	if t, ok := d.timers[process]; ok {
		t.Stop()
		delete(d.timers, process)
	}
	d.mu.Unlock()
	d.flush(process)
}

// PendingCount returns the number of buffered entries for a process.
func (d *Digest) PendingCount(process string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.entries[process])
}

func (d *Digest) flush(process string) {
	d.mu.Lock()
	batch := d.entries[process]
	delete(d.entries, process)
	delete(d.timers, process)
	d.mu.Unlock()

	if len(batch) > 0 && d.onFlush != nil {
		d.onFlush(process, batch)
	}
}
