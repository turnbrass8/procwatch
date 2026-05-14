package notify

import (
	"sync"
	"time"
)

// Level represents an escalation severity level.
type Level int

const (
	LevelNormal  Level = iota
	LevelWarning       // threshold exceeded once
	LevelCritical      // threshold exceeded repeatedly
)

// EscalationPolicy defines when to escalate alert levels.
type EscalationPolicy struct {
	WarningAfter  int           // number of consecutive threshold hits
	CriticalAfter int           // number of consecutive threshold hits
	ResetAfter    time.Duration // reset level after this quiet period
}

// DefaultEscalationPolicy returns sensible defaults.
func DefaultEscalationPolicy() EscalationPolicy {
	return EscalationPolicy{
		WarningAfter:  2,
		CriticalAfter: 5,
		ResetAfter:    5 * time.Minute,
	}
}

type processState struct {
	hits    int
	level   Level
	lastHit time.Time
}

// Escalator tracks consecutive threshold violations per process
// and returns the current alert level.
type Escalator struct {
	mu     sync.Mutex
	policy EscalationPolicy
	state  map[string]*processState
}

// NewEscalator creates an Escalator with the given policy.
func NewEscalator(policy EscalationPolicy) *Escalator {
	return &Escalator{
		policy: policy,
		state:  make(map[string]*processState),
	}
}

// Record registers a threshold violation for the named process and
// returns the resulting escalation level.
func (e *Escalator) Record(name string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	s, ok := e.state[name]
	if !ok {
		s = &processState{}
		e.state[name] = s
	}

	now := time.Now()
	if s.lastHit != (time.Time{}) && now.Sub(s.lastHit) > e.policy.ResetAfter {
		s.hits = 0
		s.level = LevelNormal
	}

	s.hits++
	s.lastHit = now

	switch {
	case s.hits >= e.policy.CriticalAfter:
		s.level = LevelCritical
	case s.hits >= e.policy.WarningAfter:
		s.level = LevelWarning
	default:
		s.level = LevelNormal
	}

	return s.level
}

// Reset clears the escalation state for a named process.
func (e *Escalator) Reset(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.state, name)
}

// CurrentLevel returns the current escalation level without recording a hit.
func (e *Escalator) CurrentLevel(name string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()
	if s, ok := e.state[name]; ok {
		return s.level
	}
	return LevelNormal
}
