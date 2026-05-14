package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/procwatch/internal/notify"
)

func newRateLimiter() *notify.RateLimiter {
	return notify.NewRateLimiter(time.Minute, 5)
}

func TestHandleRateLimitStatus_MissingName(t *testing.T) {
	h := &rateLimitHandler{rl: newRateLimiter()}
	req := httptest.NewRequest(http.MethodGet, "/ratelimit/status", nil)
	w := httptest.NewRecorder()
	h.handleRateLimitStatus(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleRateLimitStatus_Known(t *testing.T) {
	rl := newRateLimiter()
	rl.Allow("nginx")
	rl.Allow("nginx")
	h := &rateLimitHandler{rl: rl}

	req := httptest.NewRequest(http.MethodGet, "/ratelimit/status?name=nginx", nil)
	w := httptest.NewRecorder()
	h.handleRateLimitStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["process"] != "nginx" {
		t.Errorf("unexpected process: %v", resp["process"])
	}
	if resp["count"].(float64) != 2 {
		t.Errorf("expected count 2, got %v", resp["count"])
	}
}

func TestHandleRateLimitReset_Named(t *testing.T) {
	rl := newRateLimiter()
	rl.Allow("redis")
	h := &rateLimitHandler{rl: rl}

	body := bytes.NewBufferString(`{"name":"redis"}`)
	req := httptest.NewRequest(http.MethodPost, "/ratelimit/reset", body)
	w := httptest.NewRecorder()
	h.handleRateLimitReset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if rl.Count("redis") != 0 {
		t.Error("expected count to be 0 after reset")
	}
}

func TestHandleRateLimitReset_All(t *testing.T) {
	rl := newRateLimiter()
	rl.Allow("a")
	rl.Allow("b")
	h := &rateLimitHandler{rl: rl}

	req := httptest.NewRequest(http.MethodPost, "/ratelimit/reset", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	h.handleRateLimitReset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandleRateLimitReset_MethodNotAllowed(t *testing.T) {
	h := &rateLimitHandler{rl: newRateLimiter()}
	req := httptest.NewRequest(http.MethodGet, "/ratelimit/reset", nil)
	w := httptest.NewRecorder()
	h.handleRateLimitReset(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
