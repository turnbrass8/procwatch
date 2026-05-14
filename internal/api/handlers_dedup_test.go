package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"procwatch/internal/notify"
)

func newDeduplicator() *notify.Deduplicator {
	return notify.NewDeduplicator(30 * time.Second)
}

func TestHandleDedupReset_AllProcesses(t *testing.T) {
	d := newDeduplicator()
	srv := &Server{dedup: d}

	req := httptest.NewRequest(http.MethodPost, "/dedup/reset", nil)
	w := httptest.NewRecorder()
	srv.handleDedupReset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "reset" {
		t.Errorf("expected status=reset, got %q", resp["status"])
	}
}

func TestHandleDedupReset_NamedProcess(t *testing.T) {
	d := newDeduplicator()
	srv := &Server{dedup: d}

	req := httptest.NewRequest(http.MethodPost, "/dedup/reset?name=nginx", nil)
	w := httptest.NewRecorder()
	srv.handleDedupReset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["process"] != "nginx" {
		t.Errorf("expected process=nginx, got %q", resp["process"])
	}
}

func TestHandleDedupReset_MethodNotAllowed(t *testing.T) {
	d := newDeduplicator()
	srv := &Server{dedup: d}

	req := httptest.NewRequest(http.MethodGet, "/dedup/reset", nil)
	w := httptest.NewRecorder()
	srv.handleDedupReset(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleDedupCheck_NotSuppressed(t *testing.T) {
	d := newDeduplicator()
	srv := &Server{dedup: d}

	req := httptest.NewRequest(http.MethodGet, "/dedup/check?name=nginx&reason=crash", nil)
	w := httptest.NewRecorder()
	srv.handleDedupCheck(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["suppressed"].(bool) {
		t.Error("expected suppressed=false on first check")
	}
}

func TestHandleDedupCheck_MissingName(t *testing.T) {
	d := newDeduplicator()
	srv := &Server{dedup: d}

	req := httptest.NewRequest(http.MethodGet, "/dedup/check?reason=crash", nil)
	w := httptest.NewRecorder()
	srv.handleDedupCheck(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
