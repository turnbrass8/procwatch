package notify

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// DedupWindow is the duration during which identical alerts are suppressed.
const DedupWindow = 5 * time.Minute

// dedupEntry tracks when an alert fingerprint was last seen.
type dedupEntry struct {
	seenAt time.Time
}

// Deduplicator suppresses repeated identical alerts within a time window.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]dedupEntry
	now     func() time.Time
}

// NewDeduplicator creates a Deduplicator with the given suppression window.
func NewDeduplicator(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window:  window,
		entries: make(map[string]dedupEntry),
		now:     time.Now,
	}
}

// IsDuplicate returns true if an identical alert was already seen within the window.
// If not a duplicate, the fingerprint is recorded.
func (d *Deduplicator) IsDuplicate(process, reason string) bool {
	key := fingerprint(process, reason)
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if e, ok := d.entries[key]; ok {
		if now.Sub(e.seenAt) < d.window {
			return true
		}
	}
	d.entries[key] = dedupEntry{seenAt: now}
	return false
}

// Reset clears the dedup state for a specific process, or all processes if name is empty.
func (d *Deduplicator) Reset(process string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if process == "" {
		d.entries = make(map[string]dedupEntry)
		return
	}
	for k := range d.entries {
		if len(k) >= len(process) && k[:len(process)] == process {
			delete(d.entries, k)
		}
	}
}

// fingerprint produces a stable key from process name and alert reason.
func fingerprint(process, reason string) string {
	h := sha256.Sum256([]byte(process + "|" + reason))
	return fmt.Sprintf("%x", h[:8])
}
