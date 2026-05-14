package api

import (
	"net/http"
)

// handleDedupReset resets deduplication state for all or a named process.
// POST /dedup/reset?name=<optional>
func (s *Server) handleDedupReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		s.dedup.Reset("")
		writeJSON(w, map[string]string{
			"status":  "reset",
			"scope":   "all",
		})
		return
	}
	s.dedup.Reset(name)
	writeJSON(w, map[string]string{
		"status":  "reset",
		"process": name,
	})
}

// handleDedupCheck reports whether an alert for a given process+reason would be suppressed.
// GET /dedup/check?name=<name>&reason=<reason>
func (s *Server) handleDedupCheck(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing 'name' query parameter", http.StatusBadRequest)
		return
	}
	reason := r.URL.Query().Get("reason")
	if reason == "" {
		reason = "unknown"
	}

	// IsDuplicate does not record — use IsDuplicate for a read-only check.
	suppressed := s.dedup.IsDuplicate(name, reason)
	writeJSON(w, map[string]interface{}{
		"process":    name,
		"reason":     reason,
		"suppressed": suppressed,
	})
}
