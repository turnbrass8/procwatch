package history

import (
	"sync"
	"time"
)

// EventType categorizes what kind of event occurred.
type EventType string

const (
	EventCrash     EventType = "crash"
	EventThreshold EventType = "threshold"
	EventRecovered EventType = "recovered"
)

// Event represents a single recorded monitor event for a process.
type Event struct {
	ProcessName string
	Type        EventType
	Message     string
	OccurredAt  time.Time
}

// Store holds an in-memory ring buffer of recent events per process.
type Store struct {
	mu       sync.RWMutex
	events   map[string][]Event
	maxPerProcess int
}

// NewStore creates a Store that retains up to maxPerProcess events per process.
func NewStore(maxPerProcess int) *Store {
	if maxPerProcess <= 0 {
		maxPerProcess = 50
	}
	return &Store{
		events:        make(map[string][]Event),
		maxPerProcess: maxPerProcess,
	}
}

// Record appends a new event, evicting the oldest if the buffer is full.
func (s *Store) Record(e Event) {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	buf := s.events[e.ProcessName]
	if len(buf) >= s.maxPerProcess {
		buf = buf[1:]
	}
	s.events[e.ProcessName] = append(buf, e)
}

// Get returns a copy of recorded events for the given process name.
func (s *Store) Get(processName string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.events[processName]
	out := make([]Event, len(src))
	copy(out, src)
	return out
}

// All returns a flat slice of every recorded event across all processes.
func (s *Store) All() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Event
	for _, evts := range s.events {
		out = append(out, evts...)
	}
	return out
}

// Clear removes all events for a given process.
func (s *Store) Clear(processName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, processName)
}
