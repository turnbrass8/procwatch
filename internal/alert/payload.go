package alert

import "time"

// Reason describes why an alert was triggered.
type Reason string

const (
	ReasonNotRunning  Reason = "process_not_running"
	ReasonCPUExceeded Reason = "cpu_exceeded"
	ReasonMemExceeded Reason = "mem_exceeded"
)

// Payload is the JSON body sent to the webhook on an alert event.
type Payload struct {
	Process   string    `json:"process"`
	Reason    Reason    `json:"reason"`
	Message   string    `json:"message"`
	CPUPct    float64   `json:"cpu_pct,omitempty"`
	MemMB     float64   `json:"mem_mb,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Send dispatches p to the configured webhook URL.
func (s *Sender) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}
	return s.send(p)
}
