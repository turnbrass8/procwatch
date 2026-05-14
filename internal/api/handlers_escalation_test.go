package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/procwatch/internal/notify"
)

func newEscalator() *notify.Escalator {
	return notify.NewEscalator(notify.EscalationPolicy{
		WarningAfter:  2,
		CriticalAfter: 4,
		ResetAfter:    5 * 60 * 1000000000, // 5 min
	})
}

func TestHandleEscalationStatus_Normal(t *testing.T) {
	esc := newEscalator()
	h := handleEscalationStatus(esc)
	req := httptest.NewRequest(http.MethodGet, "/escalation/status?name=svc", nil)
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["level"] != "normal" {
		t.Errorf("expected normal, got %v", resp["level"])
	}
}

func TestHandleEscalationStatus_MissingName(t *testing.T) {
	esc := newEscalator()
	h := handleEscalationStatus(esc)
	req := httptest.NewRequest(http.MethodGet, "/escalation/status", nil)
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleEscalationStatus_Critical(t *testing.T) {
	esc := newEscalator()
	for i := 0; i < 4; i++ {
		esc.Record("svc")
	}
	h := handleEscalationStatus(esc)
	req := httptest.NewRequest(http.MethodGet, "/escalation/status?name=svc", nil)
	w := httptest.NewRecorder()
	h(w, req)
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["level"] != "critical" {
		t.Errorf("expected critical, got %v", resp["level"])
	}
}

func TestHandleEscalationReset_OK(t *testing.T) {
	esc := newEscalator()
	for i := 0; i < 4; i++ {
		esc.Record("svc")
	}
	h := handleEscalationReset(esc)
	req := httptest.NewRequest(http.MethodPost, "/escalation/reset?name=svc", nil)
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if esc.CurrentLevel("svc") != notify.LevelNormal {
		t.Error("expected level to be reset to Normal")
	}
}

func TestHandleEscalationReset_MethodNotAllowed(t *testing.T) {
	esc := newEscalator()
	h := handleEscalationReset(esc)
	req := httptest.NewRequest(http.MethodGet, "/escalation/reset?name=svc", nil)
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}
