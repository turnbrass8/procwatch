package notify

import (
	"testing"
	"time"
)

// TestSchedulerWithSilencer verifies that a Scheduler and Silencer work
// independently and both must allow an alert for it to be delivered.
func TestSchedulerWithSilencer(t *testing.T) {
	// Fix time to Wednesday 14:00
	now := func() time.Time {
		return time.Date(2024, 1, 3, 14, 0, 0, 0, time.UTC) // Wednesday
	}

	sched := NewScheduler(now)
	sched.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Wednesday)},
		Hours: TimeRange{Start: "09:00", End: "17:00"},
	})

	silencer := NewSilencer(now)

	// Both allow: alert should pass
	if !sched.IsActive("myapp") {
		t.Fatal("scheduler should be active")
	}
	if silencer.IsSilenced("myapp") {
		t.Fatal("silencer should not be active")
	}

	// Silence the process: alert should be suppressed
	silencer.Silence("myapp", 5*time.Minute)
	if !silencer.IsSilenced("myapp") {
		t.Fatal("expected process to be silenced")
	}

	// Schedule outside window: alert should be suppressed
	nowNight := func() time.Time {
		return time.Date(2024, 1, 3, 23, 0, 0, 0, time.UTC)
	}
	schedNight := NewScheduler(nowNight)
	schedNight.Set("myapp", ScheduleEntry{
		Days:  []DayOfWeek{DayOfWeek(time.Wednesday)},
		Hours: TimeRange{Start: "09:00", End: "17:00"},
	})
	if schedNight.IsActive("myapp") {
		t.Fatal("scheduler should be inactive at night")
	}

	// After lifting silence and within window, alert passes
	silencer.Lift("myapp")
	if silencer.IsSilenced("myapp") {
		t.Fatal("expected silence to be lifted")
	}
	if !sched.IsActive("myapp") {
		t.Fatal("scheduler should still be active in original window")
	}
}
