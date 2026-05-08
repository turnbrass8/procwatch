package monitor

import (
	"fmt"
	"log"

	"github.com/user/procwatch/internal/alert"
	"github.com/user/procwatch/internal/config"
)

// CheckProcess looks up the process by name, retrieves its stats, and sends
// webhook alerts if the process is not found or exceeds configured thresholds.
func CheckProcess(proc config.Process, sender *alert.Sender) {
	pid, err := FindPID(proc.Name)
	if err != nil {
		payload := alert.Payload{
			ProcessName: proc.Name,
			Event:       "crash",
			Message:     fmt.Sprintf("process '%s' not found: %v", proc.Name, err),
		}
		if alertErr := sender.Send(payload); alertErr != nil {
			log.Printf("[ERROR] failed to send crash alert for %s: %v", proc.Name, alertErr)
		}
		return
	}

	stats, err := GetStats(pid)
	if err != nil {
		log.Printf("[WARN] could not retrieve stats for %s (pid %d): %v", proc.Name, pid, err)
		return
	}

	if proc.MaxCPU > 0 && stats.CPUPercent > proc.MaxCPU {
		payload := alert.Payload{
			ProcessName: proc.Name,
			PID:         pid,
			Event:       "cpu_threshold",
			Message:     fmt.Sprintf("CPU %.2f%% exceeds limit %.2f%%", stats.CPUPercent, proc.MaxCPU),
			CPUPercent:  stats.CPUPercent,
			MemoryMB:    stats.MemoryMB,
		}
		if alertErr := sender.Send(payload); alertErr != nil {
			log.Printf("[ERROR] failed to send CPU alert for %s: %v", proc.Name, alertErr)
		}
	}

	if proc.MaxMemMB > 0 && stats.MemoryMB > proc.MaxMemMB {
		payload := alert.Payload{
			ProcessName: proc.Name,
			PID:         pid,
			Event:       "memory_threshold",
			Message:     fmt.Sprintf("memory %.2f MB exceeds limit %.2f MB", stats.MemoryMB, proc.MaxMemMB),
			CPUPercent:  stats.CPUPercent,
			MemoryMB:    stats.MemoryMB,
		}
		if alertErr := sender.Send(payload); alertErr != nil {
			log.Printf("[ERROR] failed to send memory alert for %s: %v", proc.Name, alertErr)
		}
	}
}
