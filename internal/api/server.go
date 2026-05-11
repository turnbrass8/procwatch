package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/procwatch/internal/history"
)

// Server exposes a lightweight HTTP API for querying process history.
type Server struct {
	store  *history.Store
	addr   string
	httpSrv *http.Server
}

// NewServer creates a new API server bound to addr.
func NewServer(addr string, store *history.Store) *Server {
	s := &Server{
		store: store,
		addr:  addr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/history", s.handleHistory)
	mux.HandleFunc("/summary", s.handleSummary)

	s.httpSrv = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s
}

// Start begins serving HTTP requests. It blocks until the server stops.
func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	return s.httpSrv.Close()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
