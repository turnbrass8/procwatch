package notify

import (
	"testing"
	"time"
)

func fixedTime(weekday time.Weekday, hour, min int) func() time.Time {
	return func() time.Time {
		base := time.Date(2024, 1, 1, hour, min, 0, 0, time.UTC)
		// Shift to desired weekday (2024-01-01 is Monday = 1)
		delta := int(weekday) - int(base.Weekday())
		return base.AddDate(0, 0, delta)
	}
}

func TestScheduler_NoScheduleAlwaysActive(t *testing.T) {
	s := NewScheduler(fixedTime(time.Saturday, 3, 0))
	if !s.IsActive("myapp") {
		t.Fatal("expected active when no schedule defined")
	}
}

func TestScheduler_InWindowIsActive(t *testing.T) {
	s := NewScheduler(fixedTime(time.Monday, 10, 30))
	s.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Monday)},
		Hours: TimeRange{Start: "09:00", End: "17:00"},
	})
	if !s.IsActive("myapp") {
		t.Fatal("expected active during business hours")
	}
}

func TestScheduler_OutsideHoursNotActive(t *testing.T) {
	s := NewScheduler(fixedTime(time.Monday, 22, 0))
	s.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Monday)},
		Hours: TimeRange{Start: "09:00", End: "17:00"},
	})
	if s.IsActive("myapp") {
		t.Fatal("expected inactive outside hours")
	}
}

func TestScheduler_WrongDayNotActive(t *testing.T) {
	s := NewScheduler(fixedTime(time.Sunday, 12, 0))
	s.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Monday), DayOfWeek(time.Friday)},
		Hours: TimeRange{Start: "00:00", End: "23:59"},
	})
	if s.IsActive("myapp") {
		t.Fatal("expected inactive on Sunday")
	}
}

func TestScheduler_GlobalFallback(t *testing.T) {
	s := NewScheduler(fixedTime(time.Tuesday, 14, 0))
	s.Set("*", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Monday), DayOfWeek(time.Tuesday)},
		Hours: TimeRange{Start: "08:00", End: "18:00"},
	})
	if !s.IsActive("unknown-service") {
		t.Fatal("expected global schedule to apply")
	}
}

func TestScheduler_ProcessOverridesGlobal(t *testing.T) {
	s := NewScheduler(fixedTime(time.Tuesday, 14, 0))
	s.Set("*", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Tuesday)},
		Hours: TimeRange{Start: "08:00", End: "18:00"},
	})
	// Process-specific schedule blocks Tuesday
	s.Set("critical", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Wednesday)},
		Hours: TimeRange{Start: "08:00", End: "18:00"},
	})
	if s.IsActive("critical") {
		t.Fatal("process-specific schedule should override global")
	}
}

func TestScheduler_RemoveClearsSchedule(t *testing.T) {
	s := NewScheduler(fixedTime(time.Sunday, 12, 0))
	s.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Monday)},
		Hours: TimeRange{Start: "09:00", End: "17:00"},
	})
	s.Remove("myapp")
	if !s.IsActive("myapp") {
		t.Fatal("expected active after schedule removed")
	}
}
