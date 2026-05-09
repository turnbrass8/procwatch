package monitor

import (
	"log"
	"time"

	"github.com/user/procwatch/internal/alert"
	"github.com/user/procwatch/internal/config"
)

// Runner periodically checks all configured processes and sends alerts on issues.
type Runner struct {
	cfg    *config.Config
	sender *alert.Sender
}

// NewRunner creates a Runner with the given config and alert sender.
func NewRunner(cfg *config.Config, sender *alert.Sender) *Runner {
	return &Runner{cfg: cfg, sender: sender}
}

// Run starts the monitoring loop. It blocks until the provided stop channel is closed.
func (r *Runner) Run(stop <-chan struct{}) {
	interval := time.Duration(r.cfg.IntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("procwatch: starting monitor loop (interval=%s, processes=%d)",
		interval, len(r.cfg.Processes))

	for {
		select {
		case <-ticker.C:
			r.checkAll()
		case <-stop:
			log.Println("procwatch: monitor loop stopped")
			return
		}
	}
}

func (r *Runner) checkAll() {
	for _, proc := range r.cfg.Processes {
		result := CheckProcess(proc)
		if result == nil {
			continue
		}
		if err := r.sender.Send(*result); err != nil {
			log.Printf("procwatch: failed to send alert for %s: %v", proc.Name, err)
		}
	}
}
