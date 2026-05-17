package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

type scheduleRequest struct {
	Process string            `json:"process"`
	Days    []int             `json:"days"`
	Start   string            `json:"start"`
	End     string            `json:"end"`
}

func handleScheduleSet(sched *notify.Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req scheduleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Process == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if req.Start == "" || req.End == "" || len(req.Days) == 0 {
			http.Error(w, "days, start, and end are required", http.StatusBadRequest)
			return
		}
		days := make([]notify.DayOfWeek, len(req.Days))
		for i, d := range req.Days {
			days[i] = notify.DayOfWeek(d)
		}
		sched.Set(req.Process, notify.ScheduleEntry{
			Days:  days,
			Hours: notify.TimeRange{Start: req.Start, End: req.End},
		})
		writeJSON(w, map[string]string{"status": "ok", "process": req.Process})
	}
}

func handleScheduleRemove(sched *notify.Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		sched.Remove(name)
		writeJSON(w, map[string]string{"status": "removed", "process": name})
	}
}

func handleScheduleStatus(sched *notify.Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		active := sched.IsActive(name)
		writeJSON(w, map[string]interface{}{"process": name, "active": active})
	}
}
