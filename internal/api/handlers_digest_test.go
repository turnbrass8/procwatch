package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"procwatch/internal/notify"
)

func newDigest() *notify.Digest {
	return notify.NewDigest(100*time.Millisecond, 10)
}

func TestHandleDigestFlush_AllProcesses(t *testing.T) {
	d := newDigest()
	d.Add("svcA", "cpu", "CPU over threshold")
	d.Add("svcB", "crash", "process crashed")

	s := newTestServer()
	s.digest = d

	req := httptest.NewRequest(http.MethodPost, "/digest/flush", nil)
	w := httptest.NewRecorder()
	s.handleDigestFlush(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["flushed"] == nil {
		t.Error("expected flushed field in response")
	}
}

func TestHandleDigestFlush_NamedProcess(t *testing.T) {
	d := newDigest()
	d.Add("svcA", "cpu", "CPU over threshold")

	s := newTestServer()
	s.digest = d

	req := httptest.NewRequest(http.MethodPost, "/digest/flush?name=svcA", nil)
	w := httptest.NewRecorder()
	s.handleDigestFlush(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandleDigestFlush_MethodNotAllowed(t *testing.T) {
	s := newTestServer()
	s.digest = newDigest()

	req := httptest.NewRequest(http.MethodGet, "/digest/flush", nil)
	w := httptest.NewRecorder()
	s.handleDigestFlush(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleDigestStatus_Empty(t *testing.T) {
	s := newTestServer()
	s.digest = newDigest()

	req := httptest.NewRequest(http.MethodGet, "/digest/status?name=svcX", nil)
	w := httptest.NewRecorder()
	s.handleDigestStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["pending"] == nil {
		t.Error("expected pending field")
	}
}

func TestHandleDigestStatus_MissingName(t *testing.T) {
	s := newTestServer()
	s.digest = newDigest()

	req := httptest.NewRequest(http.MethodGet, "/digest/status", nil)
	w := httptest.NewRecorder()
	s.handleDigestStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
