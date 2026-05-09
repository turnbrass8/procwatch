package monitor

import (
	"log"

	"github.com/user/procwatch/internal/alert"
	"github.com/user/procwatch/internal/config"
)

// CheckAndAlert evaluates a process against its configured thresholds and
// sends a webhook alert (with retries) if a violation is detected.
func CheckAndAlert(proc config.Process, sender *alert.Sender, retryCfg alert.RetryConfig) {
	result := CheckProcess(proc)
	if result == nil {
		return
	}

	p := alert.Payload{
		Process: result.Name,
		PID:     result.PID,
		Event:   result.Reason,
		CPU:     result.CPU,
		MemMB:   result.MemMB,
	}

	if err := sender.SendWithRetry(p, retryCfg); err != nil {
		log.Printf("[procwatch] failed to deliver alert for %q after retries: %v", proc.Name, err)
	}
}
