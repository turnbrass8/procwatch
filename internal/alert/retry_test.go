package alert_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/procwatch/internal/alert"
)

func TestSendWithRetry_SucceedsOnFirstAttempt(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := alert.NewSender(ts.URL)
	cfg := alert.RetryConfig{MaxAttempts: 3, Delay: 10 * time.Millisecond}
	err := s.SendWithRetry(alert.Payload{Process: "svc", Event: "crash"}, cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestSendWithRetry_RetriesOnFailure(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := alert.NewSender(ts.URL)
	cfg := alert.RetryConfig{MaxAttempts: 3, Delay: 10 * time.Millisecond}
	err := s.SendWithRetry(alert.Payload{Process: "svc", Event: "crash"}, cfg)
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestSendWithRetry_ExhaustsAttempts(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	s := alert.NewSender(ts.URL)
	cfg := alert.RetryConfig{MaxAttempts: 3, Delay: 10 * time.Millisecond}
	err := s.SendWithRetry(alert.Payload{Process: "svc", Event: "crash"}, cfg)
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := alert.DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 2*time.Second {
		t.Errorf("expected Delay=2s, got %v", cfg.Delay)
	}
}
