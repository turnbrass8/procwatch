package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the webhook alert payload sent when a process event occurs.
type Payload struct {
	Timestamp   time.Time `json:"timestamp"`
	ProcessName string    `json:"process_name"`
	PID         int32     `json:"pid"`
	Event       string    `json:"event"`
	Message     string    `json:"message"`
	CPUPercent  float64   `json:"cpu_percent,omitempty"`
	MemoryMB    float32   `json:"memory_mb,omitempty"`
}

// Sender handles sending alert payloads to a webhook URL.
type Sender struct {
	WebhookURL string
	Client     *http.Client
}

// NewSender creates a new Sender with a default HTTP client timeout.
func NewSender(webhookURL string) *Sender {
	return &Sender{
		WebhookURL: webhookURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send marshals the payload to JSON and POSTs it to the configured webhook URL.
func (s *Sender) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("alert: failed to marshal payload: %w", err)
	}

	resp, err := s.Client.Post(s.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}
