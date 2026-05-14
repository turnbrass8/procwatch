package api

import (
	"encoding/json"
	"net/http"

	"procwatch/internal/notify"
)

type filterHandlers struct {
	filter *notify.Filter
}

func (h *filterHandlers) handleFilterAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Process  string   `json:"process"`
		Keywords []string `json:"keywords"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Process == "" || len(body.Keywords) == 0 {
		http.Error(w, "invalid request: process and keywords required", http.StatusBadRequest)
		return
	}
	h.filter.AddRule(body.Process, body.Keywords)
	writeJSON(w, map[string]string{"status": "added", "process": body.Process})
}

func (h *filterHandlers) handleFilterRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}
	h.filter.RemoveRule(name)
	writeJSON(w, map[string]string{"status": "removed", "process": name})
}

func (h *filterHandlers) handleFilterStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rules := h.filter.Rules()
	writeJSON(w, map[string]interface{}{"rules": rules})
}
