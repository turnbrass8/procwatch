package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/procwatch/internal/notify"
)

func newSilencer() *notify.Silencer {
	return notify.NewSilencer()
}

func TestHandleSilenceCreate_OK(t *testing.T) {
	s := newSilencer()
	body := `{"process":"nginx","duration":"10m"}`
	req := httptest.NewRequest(http.MethodPost, "/silence", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceCreate(s)(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !s.IsSilenced("nginx") {
		t.Fatal("expected nginx to be silenced")
	}
}

func TestHandleSilenceCreate_InvalidDuration(t *testing.T) {
	s := newSilencer()
	body := `{"process":"nginx","duration":"bad"}`
	req := httptest.NewRequest(http.MethodPost, "/silence", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceCreate(s)(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleSilenceCreate_MethodNotAllowed(t *testing.T) {
	s := newSilencer()
	req := httptest.NewRequest(http.MethodGet, "/silence", nil)
	w := httptest.NewRecorder()
	handleSilenceCreate(s)(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestHandleSilenceLift_Named(t *testing.T) {
	s := newSilencer()
	body := `{"process":"nginx","duration":"10m"}`
	req := httptest.NewRequest(http.MethodPost, "/silence", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleSilenceCreate(s)(w, req)

	req2 := httptest.NewRequest(http.MethodPost, "/silence/lift?name=nginx", nil)
	w2 := httptest.NewRecorder()
	handleSilenceLift(s)(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	if s.IsSilenced("nginx") {
		t.Fatal("expected silence to be lifted")
	}
}

func TestHandleSilenceStatus_NotSilenced(t *testing.T) {
	s := newSilencer()
	req := httptest.NewRequest(http.MethodGet, "/silence/status?name=nginx", nil)
	w := httptest.NewRecorder()
	handleSilenceStatus(s)(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["silenced"].(bool) {
		t.Fatal("expected not silenced")
	}
}

func TestHandleSilenceStatus_MissingName(t *testing.T) {
	s := newSilencer()
	req := httptest.NewRequest(http.MethodGet, "/silence/status", nil)
	w := httptest.NewRecorder()
	handleSilenceStatus(s)(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "missing name") {
		t.Fatal("expected missing name error")
	}
}
