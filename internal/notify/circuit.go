package notify

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota // normal operation
	StateOpen                // blocking alerts
	StateHalfOpen            // testing recovery
)

// CircuitBreaker stops alerting for a process after too many consecutive
// failures, then probes periodically to see if the endpoint has recovered.
type CircuitBreaker struct {
	mu           sync.Mutex
	states       map[string]*circuitState
	maxFailures  int
	openDuration time.Duration
}

type circuitState struct {
	state       State
	failures    int
	openedAt    time.Time
}

// NewCircuitBreaker returns a CircuitBreaker that opens after maxFailures
// consecutive failures and resets after openDuration.
func NewCircuitBreaker(maxFailures int, openDuration time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		states:       make(map[string]*circuitState),
		maxFailures:  maxFailures,
		openDuration: openDuration,
	}
}

func (cb *CircuitBreaker) getOrCreate(process string) *circuitState {
	if s, ok := cb.states[process]; ok {
		return s
	}
	s := &circuitState{state: StateClosed}
	cb.states[process] = s
	return s
}

// Allow returns true if an alert should be sent for the given process.
func (cb *CircuitBreaker) Allow(process string) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s := cb.getOrCreate(process)
	switch s.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(s.openedAt) >= cb.openDuration {
			s.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful alert delivery, closing the circuit.
func (cb *CircuitBreaker) RecordSuccess(process string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s := cb.getOrCreate(process)
	s.failures = 0
	s.state = StateClosed
}

// RecordFailure records a failed alert delivery. After maxFailures the
// circuit opens and alerts are suppressed for openDuration.
func (cb *CircuitBreaker) RecordFailure(process string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	s := cb.getOrCreate(process)
	s.failures++
	if s.failures >= cb.maxFailures {
		s.state = StateOpen
		s.openedAt = time.Now()
	}
}

// StateFor returns the current State for a process.
func (cb *CircuitBreaker) StateFor(process string) State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.getOrCreate(process).state
}

// Reset clears circuit state for a specific process (or all if empty).
func (cb *CircuitBreaker) Reset(process string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if process == "" {
		cb.states = make(map[string]*circuitState)
		return
	}
	delete(cb.states, process)
}
