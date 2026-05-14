package api

import (
	"net/http"
)

// handleDigestFlush forces a flush of pending digest alerts.
// POST /digest/flush?name=<process>  (name optional; omit to flush all)
func (s *Server) handleDigestFlush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")

	var flushed map[string][]string
	if name != "" {
		batch := s.digest.Flush(name)
		flushed = map[string][]string{name: batch}
	} else {
		flushed = s.digest.FlushAll()
	}

	total := 0
	for _, msgs := range flushed {
		total += len(msgs)
	}

	writeJSON(w, map[string]interface{}{
		"flushed":   flushed,
		"total":     total,
		"processes": len(flushed),
	})
}

// handleDigestStatus returns the number of pending digest alerts for a process.
// GET /digest/status?name=<process>
func (s *Server) handleDigestStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing required query param: name", http.StatusBadRequest)
		return
	}

	count := s.digest.Pending(name)

	writeJSON(w, map[string]interface{}{
		"process": name,
		"pending": count,
	})
}
