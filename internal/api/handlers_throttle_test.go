package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/procwatch/internal/notify"
)

func newThrottle() *notify.Throttle {
	return notify.NewThrottle(5 * time.Minute)
}

func TestHandleThrottleReset_AllProcesses(t *testing.T) {
	th := newThrottle()
	th.Allow("svcA") // consume slot

	req := httptest.NewRequest(http.MethodDelete, "/throttle", nil)
	w := httptest.NewRecorder()
	handleThrottleReset(th)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "all throttles reset" {
		t.Errorf("unexpected status: %s", resp["status"])
	}
}

func TestHandleThrottleReset_NamedProcess(t *testing.T) {
	th := newThrottle()
	req := httptest.NewRequest(http.MethodDelete, "/throttle?name=myapp", nil)
	w := httptest.NewRecorder()
	handleThrottleReset(th)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["process"] != "myapp" {
		t.Errorf("expected process=myapp, got %s", resp["process"])
	}
}

func TestHandleThrottleReset_MethodNotAllowed(t *testing.T) {
	th := newThrottle()
	req := httptest.NewRequest(http.MethodGet, "/throttle", nil)
	w := httptest.NewRecorder()
	handleThrottleReset(th)(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleThrottleStatus_NotThrottled(t *testing.T) {
	th := newThrottle()
	req := httptest.NewRequest(http.MethodGet, "/throttle/status?name=freshsvc", nil)
	w := httptest.NewRecorder()
	handleThrottleStatus(th)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["throttled"].(bool) {
		t.Error("expected process to not be throttled")
	}
}

func TestHandleThrottleStatus_MissingName(t *testing.T) {
	th := newThrottle()
	req := httptest.NewRequest(http.MethodGet, "/throttle/status", nil)
	w := httptest.NewRecorder()
	handleThrottleStatus(th)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
