package history

import (
	"fmt"

	"github.com/user/procwatch/internal/monitor"
)

// Recorder wraps a Store and converts monitor check results into Events.
type Recorder struct {
	store *Store
}

// NewRecorder creates a Recorder backed by the given Store.
func NewRecorder(store *Store) *Recorder {
	return &Recorder{store: store}
}

// RecordCrash records a crash event for a process.
func (r *Recorder) RecordCrash(processName string) {
	r.store.Record(Event{
		ProcessName: processName,
		Type:        EventCrash,
		Message:     fmt.Sprintf("process %q not found or crashed", processName),
	})
}

// RecordThreshold records a resource threshold breach using process stats.
func (r *Recorder) RecordThreshold(processName string, stats monitor.ProcessStats, cpuLimit float64, memLimitMB uint64) {
	var msg string
	switch {
	case stats.CPUPercent > cpuLimit && stats.MemoryMB > memLimitMB:
		msg = fmt.Sprintf("CPU %.1f%% (limit %.1f%%) and memory %dMB (limit %dMB) exceeded",
			stats.CPUPercent, cpuLimit, stats.MemoryMB, memLimitMB)
	case stats.CPUPercent > cpuLimit:
		msg = fmt.Sprintf("CPU %.1f%% exceeded limit %.1f%%", stats.CPUPercent, cpuLimit)
	default:
		msg = fmt.Sprintf("memory %dMB exceeded limit %dMB", stats.MemoryMB, memLimitMB)
	}
	r.store.Record(Event{
		ProcessName: processName,
		Type:        EventThreshold,
		Message:     msg,
	})
}

// RecordRecovered records that a process has returned to normal thresholds.
func (r *Recorder) RecordRecovered(processName string) {
	r.store.Record(Event{
		ProcessName: processName,
		Type:        EventRecovered,
		Message:     fmt.Sprintf("process %q recovered to normal resource usage", processName),
	})
}

// RecentEvents returns the last n events for a process from the backing store.
func (r *Recorder) RecentEvents(processName string) []Event {
	return r.store.Get(processName)
}
