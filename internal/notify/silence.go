package notify

import (
	"sync"
	"time"
)

// SilenceWindow represents a named time window during which alerts are suppressed.
type SilenceWindow struct {
	Process string
	Until   time.Time
}

// Silencer suppresses alerts for specific processes during configured windows.
type Silencer struct {
	mu      sync.RWMutex
	windows map[string]time.Time
	now     func() time.Time
}

// NewSilencer creates a new Silencer.
func NewSilencer() *Silencer {
	return &Silencer{
		windows: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Silence suppresses alerts for the given process until the specified time.
func (s *Silencer) Silence(process string, until time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[process] = until
}

// IsSilenced reports whether alerts for the given process are currently suppressed.
func (s *Silencer) IsSilenced(process string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	until, ok := s.windows[process]
	if !ok {
		return false
	}
	return s.now().Before(until)
}

// Lift removes the silence window for the given process.
func (s *Silencer) Lift(process string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.windows, process)
}

// LiftAll removes all active silence windows.
func (s *Silencer) LiftAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows = make(map[string]time.Time)
}

// Status returns the silence window for a process, or zero time if not silenced.
func (s *Silencer) Status(process string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	until, ok := s.windows[process]
	if !ok || s.now().After(until) {
		return time.Time{}, false
	}
	return until, true
}
