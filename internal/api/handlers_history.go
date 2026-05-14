package api

import (
	"net/http"

	"procwatch/internal/history"
)

// handleExportJSON streams the full event history as a JSON download.
func (s *Server) handleExportJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="history.json"`)
	if err := history.ExportJSON(s.store, w); err != nil {
		http.Error(w, "export failed", http.StatusInternalServerError)
	}
}

// handleExportText streams the full event history as a plain-text download.
func (s *Server) handleExportText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="history.txt"`)
	if err := history.ExportText(s.store, w); err != nil {
		http.Error(w, "export failed", http.StatusInternalServerError)
	}
}

// handleSummarizeProcess returns a JSON summary for a single named process.
func (s *Server) handleSummarizeProcess(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing 'name' query parameter", http.StatusBadRequest)
		return
	}
	records := s.store.Get(name)
	summary := history.Summarize(name, records)
	writeJSON(w, http.StatusOK, summary)
}
