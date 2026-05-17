package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

type backoffHandlers struct {
	tracker *notify.BackoffTracker
}

func newBackoffHandlers(t *notify.BackoffTracker) *backoffHandlers {
	return &backoffHandlers{tracker: t}
}

// handleBackoffStatus returns the current attempt count and next delay for a process.
func (h *backoffHandlers) handleBackoffStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}
	attempts := h.tracker.Attempts(name)
	writeJSON(w, map[string]interface{}{
		"process":  name,
		"attempts": attempts,
	})
}

// handleBackoffReset resets the backoff counter for one or all processes.
func (h *backoffHandlers) handleBackoffReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		h.tracker.ResetAll()
		writeJSON(w, map[string]string{"status": "reset all"})
		return
	}
	h.tracker.Reset(name)
	writeJSON(w, map[string]string{"status": "reset", "process": name})
}
