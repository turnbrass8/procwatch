package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

type rateLimitHandler struct {
	rl *notify.RateLimiter
}

func (h *rateLimitHandler) handleRateLimitStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}

	count := h.rl.Count(name)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"process": name,
		"count":   count,
	})
}

func (h *rateLimitHandler) handleRateLimitReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	h.rl.Reset(body.Name)

	msg := "all processes reset"
	if body.Name != "" {
		msg = body.Name + " reset"
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": msg})
}
