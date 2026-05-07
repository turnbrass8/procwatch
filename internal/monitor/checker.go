package monitor

import (
	"fmt"

	"github.com/user/procwatch/internal/config"
)

// AlertType describes the reason an alert was triggered.
type AlertType string

const (
	AlertCrashed     AlertType = "crashed"
	AlertHighCPU     AlertType = "high_cpu"
	AlertHighMemory  AlertType = "high_memory"
)

// Alert represents a triggered threshold or crash event.
type Alert struct {
	Process string
	Type    AlertType
	Message string
}

// CheckProcess evaluates a process against its configured thresholds
// and returns any alerts that should be fired.
func CheckProcess(proc config.Process) ([]Alert, error) {
	stats, err := GetStats(proc.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for %q: %w", proc.Name, err)
	}

	var alerts []Alert

	if !stats.Alive {
		alerts = append(alerts, Alert{
			Process: proc.Name,
			Type:    AlertCrashed,
			Message: fmt.Sprintf("process %q is not running", proc.Name),
		})
		return alerts, nil
	}

	if proc.MaxCPU > 0 && stats.CPU > proc.MaxCPU {
		alerts = append(alerts, Alert{
			Process: proc.Name,
			Type:    AlertHighCPU,
			Message: fmt.Sprintf(
				"process %q CPU %.1f%% exceeds threshold %.1f%%",
				proc.Name, stats.CPU, proc.MaxCPU,
			),
		})
	}

	if proc.MaxMemoryMB > 0 && stats.Memory > proc.MaxMemoryMB {
		alerts = append(alerts, Alert{
			Process: proc.Name,
			Type:    AlertHighMemory,
			Message: fmt.Sprintf(
				"process %q memory %.1fMB exceeds threshold %.1fMB",
				proc.Name, stats.Memory, proc.MaxMemoryMB,
			),
		})
	}

	return alerts, nil
}
