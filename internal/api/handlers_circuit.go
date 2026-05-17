package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

type circuitHandlers struct {
	cb *notify.CircuitBreaker
}

func (h *circuitHandlers) handleCircuitStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name", http.StatusBadRequest)
		return
	}

	stateName := func(s notify.State) string {
		switch s {
		case notify.StateClosed:
			return "closed"
		case notify.StateOpen:
			return "open"
		case notify.StateHalfOpen:
			return "half_open"
		default:
			return "unknown"
		}
	}

	state := h.cb.StateFor(name)
	writeJSON(w, http.StatusOK, map[string]string{
		"process": name,
		"state":   stateName(state),
	})
}

func (h *circuitHandlers) handleCircuitReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	h.cb.Reset(name)

	msg := "all circuits reset"
	if name != "" {
		msg = "circuit reset for " + name
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": msg})
}
