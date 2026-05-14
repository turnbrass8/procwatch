package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

// handleEscalationStatus returns the current escalation level for a process.
func handleEscalationStatus(esc *notify.Escalator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query param", http.StatusBadRequest)
			return
		}
		level := esc.CurrentLevel(name)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"process": name,
			"level":   levelName(level),
		})
	}
}

// handleEscalationReset resets the escalation state for a process.
func handleEscalationReset(esc *notify.Escalator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query param", http.StatusBadRequest)
			return
		}
		esc.Reset(name)
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "reset",
			"process": name,
		})
	}
}

func levelName(l notify.Level) string {
	switch l {
	case notify.LevelWarning:
		return "warning"
	case notify.LevelCritical:
		return "critical"
	default:
		return "normal"
	}
}
