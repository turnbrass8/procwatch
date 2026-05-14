package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/procwatch/internal/notify"
)

type silenceRequest struct {
	Process  string `json:"process"`
	Duration string `json:"duration"`
}

func handleSilenceCreate(s *notify.Silencer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req silenceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Process == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		dur, err := time.ParseDuration(req.Duration)
		if err != nil || dur <= 0 {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		s.Silence(req.Process, time.Now().Add(dur))
		writeJSON(w, map[string]string{"status": "silenced", "process": req.Process, "duration": req.Duration})
	}
}

func handleSilenceLift(s *notify.Silencer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			s.LiftAll()
			writeJSON(w, map[string]string{"status": "all silences lifted"})
			return
		}
		s.Lift(name)
		writeJSON(w, map[string]string{"status": "lifted", "process": name})
	}
}

func handleSilenceStatus(s *notify.Silencer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}
		until, silenced := s.Status(name)
		writeJSON(w, map[string]interface{}{
			"process":  name,
			"silenced": silenced,
			"until":    until,
		})
	}
}
