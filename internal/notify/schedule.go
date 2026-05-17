package notify

import (
	"sync"
	"time"
)

// DayOfWeek represents a day (0=Sunday ... 6=Saturday).
type DayOfWeek int

// TimeRange defines a start/end time-of-day window (24h format).
type TimeRange struct {
	Start string // "HH:MM"
	End   string // "HH:MM"
}

// ScheduleEntry defines which days and time range alerts are allowed.
type ScheduleEntry struct {
	Days  []DayOfWeek
	Hours TimeRange
}

// Scheduler suppresses alerts outside of defined active windows.
type Scheduler struct {
	mu      sync.RWMutex
	schedules map[string]ScheduleEntry // keyed by process name; "*" = global
	now     func() time.Time
}

// NewScheduler creates a Scheduler with an optional clock override.
func NewScheduler(now func() time.Time) *Scheduler {
	if now == nil {
		now = time.Now
	}
	return &Scheduler{
		schedules: make(map[string]ScheduleEntry),
		now:       now,
	}
}

// Set registers an active-hours schedule for a process (use "*" for global).
func (s *Scheduler) Set(process string, entry ScheduleEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.schedules[process] = entry
}

// Remove deletes the schedule for a process.
func (s *Scheduler) Remove(process string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.schedules, process)
}

// IsActive returns true if the current time falls within the schedule for the
// given process. If no schedule is found for the process, the global "*"
// schedule is checked. If neither exists, alerts are always active.
func (s *Scheduler) IsActive(process string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.schedules[process]
	if !ok {
		entry, ok = s.schedules["*"]
		if !ok {
			return true
		}
	}
	return s.inWindow(entry)
}

func (s *Scheduler) inWindow(entry ScheduleEntry) bool {
	now := s.now()
	currentDay := DayOfWeek(now.Weekday())
	dayOK := false
	for _, d := range entry.Days {
		if d == currentDay {
			dayOK = true
			break
		}
	}
	if !dayOK {
		return false
	}

	start, err1 := time.Parse("15:04", entry.Hours.Start)
	end, err2 := time.Parse("15:04", entry.Hours.End)
	if err1 != nil || err2 != nil {
		return true
	}

	current, _ := time.Parse("15:04", now.Format("15:04"))
	return !current.Before(start) && !current.After(end)
}
