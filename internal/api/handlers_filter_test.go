package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"procwatch/internal/notify"
)

func newFilterHandlers() *filterHandlers {
	return &filterHandlers{filter: notify.NewFilter()}
}

func TestHandleFilterAdd_OK(t *testing.T) {
	h := newFilterHandlers()
	body := `{"process":"nginx","keywords":["oom","timeout"]}`
	req := httptest.NewRequest(http.MethodPost, "/filter/add", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.handleFilterAdd(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !h.filter.IsSuppressed("nginx", "oom killer hit") {
		t.Fatal("expected rule to be active after add")
	}
}

func TestHandleFilterAdd_MissingFields(t *testing.T) {
	h := newFilterHandlers()
	body := `{"process":"nginx"}`
	req := httptest.NewRequest(http.MethodPost, "/filter/add", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.handleFilterAdd(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleFilterAdd_MethodNotAllowed(t *testing.T) {
	h := newFilterHandlers()
	req := httptest.NewRequest(http.MethodGet, "/filter/add", nil)
	w := httptest.NewRecorder()
	h.handleFilterAdd(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleFilterRemove_OK(t *testing.T) {
	h := newFilterHandlers()
	h.filter.AddRule("redis", []string{"crash"})
	req := httptest.NewRequest(http.MethodDelete, "/filter/remove?name=redis", nil)
	w := httptest.NewRecorder()
	h.handleFilterRemove(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if h.filter.IsSuppressed("redis", "crash") {
		t.Fatal("expected rule to be removed")
	}
}

func TestHandleFilterRemove_MissingName(t *testing.T) {
	h := newFilterHandlers()
	req := httptest.NewRequest(http.MethodDelete, "/filter/remove", nil)
	w := httptest.NewRecorder()
	h.handleFilterRemove(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleFilterStatus_ReturnsRules(t *testing.T) {
	h := newFilterHandlers()
	h.filter.AddRule("svc", []string{"panic"})
	req := httptest.NewRequest(http.MethodGet, "/filter/status", nil)
	w := httptest.NewRecorder()
	h.handleFilterStatus(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	rules, ok := resp["rules"]
	if !ok {
		t.Fatal("expected 'rules' key in response")
	}
	rulesMap, ok := rules.(map[string]interface{})
	if !ok || rulesMap["svc"] == nil {
		t.Fatal("expected svc entry in rules")
	}
}
