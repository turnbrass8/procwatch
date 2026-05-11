package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/procwatch/internal/api"
	"github.com/user/procwatch/internal/history"
)

func newTestServer(t *testing.T) (*api.Server, *history.Store) {
	t.Helper()
	store := history.NewStore(10)
	srv := api.NewServer(":0", store)
	return srv, store
}

func TestHandleHealth(t *testing.T) {
	srv, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %q", resp["status"])
	}
}

func TestHandleHistory_MissingParam(t *testing.T) {
	srv, _ := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleHistory_KnownProcess(t *testing.T) {
	srv, store := newTestServer(t)
	store.Record("nginx", history.Record{
		Process:   "nginx",
		Timestamp: time.Now(),
		Kind:      "crash",
	})

	req := httptest.NewRequest(http.MethodGet, "/history?process=nginx", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["process"] != "nginx" {
		t.Errorf("unexpected process field: %v", resp["process"])
	}
}

func TestHandleSummary_All(t *testing.T) {
	srv, store := newTestServer(t)
	store.Record("redis", history.Record{
		Process:   "redis",
		Timestamp: time.Now(),
		Kind:      "threshold",
	})

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
