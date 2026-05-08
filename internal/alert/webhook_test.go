package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSend_Success(t *testing.T) {
	var received Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	sender := NewSender(ts.URL)
	p := Payload{
		ProcessName: "nginx",
		PID:         1234,
		Event:       "crash",
		Message:     "process exited unexpectedly",
	}

	if err := sender.Send(p); err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	if received.ProcessName != "nginx" {
		t.Errorf("expected process_name=nginx, got %s", received.ProcessName)
	}
	if received.Event != "crash" {
		t.Errorf("expected event=crash, got %s", received.Event)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sender := NewSender(ts.URL)
	err := sender.Send(Payload{Event: "crash"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_InvalidURL(t *testing.T) {
	sender := NewSender("http://127.0.0.1:0/no-listener")
	err := sender.Send(Payload{Event: "crash"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

func TestSend_TimestampAutoSet(t *testing.T) {
	before := time.Now().UTC()
	var received Payload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	sender := NewSender(ts.URL)
	sender.Send(Payload{Event: "threshold"})

	if received.Timestamp.Before(before) {
		t.Error("timestamp should be set to current time when zero")
	}
}
