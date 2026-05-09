package alert

import (
	"log"
	"time"
)

// RetryConfig holds configuration for retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
	}
}

// SendWithRetry attempts to send an alert payload, retrying on failure.
// It returns the last error encountered if all attempts fail.
func (s *Sender) SendWithRetry(p Payload, cfg RetryConfig) error {
	var err error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = s.Send(p)
		if err == nil {
			return nil
		}
		log.Printf("[procwatch] alert attempt %d/%d failed: %v", attempt, cfg.MaxAttempts, err)
		if attempt < cfg.MaxAttempts {
			time.Sleep(cfg.Delay)
		}
	}
	return err
}
