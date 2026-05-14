package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"procwatch/internal/history"
)

func seedStore(t *testing.T) *history.Store {
	t.Helper()
	st := history.NewStore(10)
	st.Record("svc", history.Record{
		Process:   "svc",
		Kind:      "crash",
		Timestamp: time.Now(),
	})
	return st
}

func TestHandleExportJSON_OK(t *testing.T) {
	srv := newTestServer(t)
	srv.store = seedStore(t)

	req := httptest.NewRequest(http.MethodGet, "/export/json", nil)
	rec := httptest.NewRecorder()
	srv.handleExportJSON(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("unexpected content-type: %s", ct)
	}
	var payload interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Errorf("body is not valid JSON: %v", err)
	}
}

func TestHandleExportText_OK(t *testing.T) {
	srv := newTestServer(t)
	srv.store = seedStore(t)

	req := httptest.NewRequest(http.MethodGet, "/export/text", nil)
	rec := httptest.NewRecorder()
	srv.handleExportText(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "svc") {
		t.Errorf("expected process name in text output, got: %s", body)
	}
}

func TestHandleSummarizeProcess_MissingName(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/summary/process", nil)
	rec := httptest.NewRecorder()
	srv.handleSummarizeProcess(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleSummarizeProcess_Known(t *testing.T) {
	srv := newTestServer(t)
	srv.store = seedStore(t)

	req := httptest.NewRequest(http.MethodGet, "/summary/process?name=svc", nil)
	rec := httptest.NewRecorder()
	srv.handleSummarizeProcess(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["process"] != "svc" {
		t.Errorf("expected process=svc, got %v", out["process"])
	}
}

func TestHandleSummarizeProcess_Unknown(t *testing.T) {
	srv := newTestServer(t)
	srv.store = history.NewStore(10)

	req := httptest.NewRequest(http.MethodGet, "/summary/process?name=unknown", nil)
	rec := httptest.NewRecorder()
	srv.handleSummarizeProcess(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
