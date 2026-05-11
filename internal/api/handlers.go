package api

import (
	"net/http"
	"time"

	"github.com/user/procwatch/internal/history"
)

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("process")
	if name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "process query parameter is required",
		})
		return
	}

	records := s.store.Get(name)
	writeJSON(w, http.StatusOK, map[string]any{
		"process": name,
		"records": records,
		"count":   len(records),
	})
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("process")
	if name != "" {
		records := s.store.Get(name)
		writeJSON(w, http.StatusOK, history.Summarize(name, records))
		return
	}

	all := s.store.All()
	summaries := history.SummarizeAll(all)
	writeJSON(w, http.StatusOK, summaries)
}
