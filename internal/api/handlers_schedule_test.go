package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/procwatch/internal/notify"
)

func newScheduler(now func() time.Time) *notify.Scheduler {
	return notify.NewScheduler(now)
}

func TestHandleScheduleSet_OK(t *testing.T) {
	sched := newScheduler(nil)
	body := `{"process":"myapp","days":[1,2,3,4,5],"start":"09:00","end":"17:00"}`
	req := httptest.NewRequest(http.MethodPost, "/schedule/set", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	handleScheduleSet(sched)(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	json.NewDecoder(rw.Body).Decode(&resp)
	if resp["process"] != "myapp" {
		t.Errorf("unexpected process: %s", resp["process"])
	}
}

func TestHandleScheduleSet_MissingFields(t *testing.T) {
	sched := newScheduler(nil)
	body := `{"process":"myapp"}`
	req := httptest.NewRequest(http.MethodPost, "/schedule/set", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	handleScheduleSet(sched)(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleScheduleSet_MethodNotAllowed(t *testing.T) {
	sched := newScheduler(nil)
	req := httptest.NewRequest(http.MethodGet, "/schedule/set", nil)
	rw := httptest.NewRecorder()
	handleScheduleSet(sched)(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}

func TestHandleScheduleRemove_OK(t *testing.T) {
	sched := newScheduler(nil)
	sched.Set("myapp", notify.ScheduleEntry{
		Days:  []notify.DayOfWeek{1},
		Hours: notify.TimeRange{Start: "09:00", End: "17:00"},
	})
	req := httptest.NewRequest(http.MethodPost, "/schedule/remove?name=myapp", nil)
	rw := httptest.NewRecorder()
	handleScheduleRemove(sched)(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestHandleScheduleRemove_MissingName(t *testing.T) {
	sched := newScheduler(nil)
	req := httptest.NewRequest(http.MethodPost, "/schedule/remove", nil)
	rw := httptest.NewRecorder()
	handleScheduleRemove(sched)(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleScheduleStatus_Active(t *testing.T) {
	// Tuesday 10:00 — should be active
	now := func() time.Time {
		return time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC) // Tuesday
	}
	sched := newScheduler(now)
	sched.Set("myapp", notify.ScheduleEntry{
		Days:  []notify.DayOfWeek{notify.DayOfWeek(time.Tuesday)},
		Hours: notify.TimeRange{Start: "09:00", End: "17:00"},
	})
	req := httptest.NewRequest(http.MethodGet, "/schedule/status?name=myapp", nil)
	rw := httptest.NewRecorder()
	handleScheduleStatus(sched)(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(rw.Body).Decode(&resp)
	if resp["active"] != true {
		t.Errorf("expected active=true, got %v", resp["active"])
	}
}

func TestHandleScheduleStatus_MissingName(t *testing.T) {
	sched := newScheduler(nil)
	req := httptest.NewRequest(http.MethodGet, "/schedule/status", nil)
	rw := httptest.NewRecorder()
	handleScheduleStatus(sched)(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}
