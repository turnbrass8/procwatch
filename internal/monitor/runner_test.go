package monitor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/procwatch/internal/alert"
	"github.com/user/procwatch/internal/config"
)

func newTestRunner(t *testing.T, webhookURL string) *Runner {
	t.Helper()
	cfg := &config.Config{
		WebhookURL:      webhookURL,
		IntervalSeconds: 1,
		Processes: []config.ProcessConfig{
			{Name: "definitely-not-running-xyz", MaxCPU: 90, MaxMemMB: 512},
		},
	}
	sender := alert.NewSender(webhookURL, 3*time.Second)
	return NewRunner(cfg, sender)
}

func TestRunner_StopsCleanly(t *testing.T) {
	var called atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	r := newTestRunner(t, ts.URL)
	stop := make(chan struct{})

	done := make(chan struct{})
	go func() {
		r.Run(stop)
		close(done)
	}()

	// Let at least one tick fire (interval=1s, wait 1.5s)
	time.Sleep(1500 * time.Millisecond)
	close(stop)

	select {
	case <-done:
		// good
	case <-time.After(2 * time.Second):
		t.Fatal("runner did not stop within timeout")
	}

	// The process doesn't exist, so an alert should have been fired at least once.
	if called.Load() == 0 {
		t.Error("expected at least one webhook call for missing process")
	}
}

func TestRunner_AlertPayloadValid(t *testing.T) {
	var payload alert.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("invalid JSON payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	r := newTestRunner(t, ts.URL)
	stop := make(chan struct{})
	go r.Run(stop)
	time.Sleep(1500 * time.Millisecond)
	close(stop)

	time.Sleep(200 * time.Millisecond) // let last send complete
	if payload.Process == "" {
		t.Error("expected non-empty process name in alert payload")
	}
}
