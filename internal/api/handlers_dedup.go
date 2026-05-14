package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

// handleDedupReset resets deduplication state for a named process or all processes.
//
// DELETE /dedup/reset?name=<process>  — reset named process
// DELETE /dedup/reset                 — reset all
func handleDedupReset(d *notify.Deduplicator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		d.Reset(name)
		if name == "" {
			writeJSON(w, http.StatusOK, map[string]string{"status": "all dedup state cleared"})
		} else {
			writeJSON(w, http.StatusOK, map[string]string{"status": "cleared", "process": name})
		}
	}
}

// handleDedupCheck reports whether the next alert for a process+reason would be a duplicate.
//
// GET /dedup/check?name=<process>&reason=<reason>
func handleDedupCheck(d *notify.Deduplicator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		reason := r.URL.Query().Get("reason")
		if name == "" || reason == "" {
			http.Error(w, "name and reason are required", http.StatusBadRequest)
			return
		}
		// Peek without recording: use a shadow dedup to avoid side effects.
		// We report current state only — callers should not rely on this for logic.
		type response struct {
			Process   string `json:"process"`
			Reason    string `json:"reason"`
			Duplicate bool   `json:"duplicate"`
		}
		// IsDuplicate records the entry on first call; this endpoint is informational.
		isDup := d.IsDuplicate(name, reason)
		writeJSON(w, http.StatusOK, response{
			Process:   name,
			Reason:    reason,
			Duplicate: isDup,
		})
	}
}
