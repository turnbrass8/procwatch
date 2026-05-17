package api

import (
	"encoding/json"
	"net/http"

	"procwatch/internal/notify"
)

type templateHandlers struct {
	reg *notify.TemplateRegistry
}

func newTemplateHandlers(reg *notify.TemplateRegistry) *templateHandlers {
	return &templateHandlers{reg: reg}
}

// handleTemplateRegister registers or updates a named webhook template.
// POST /webhooks/templates  body: {"name":"...","url":"...","headers":{...}}
func (h *templateHandlers) handleTemplateRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Name    string            `json:"name"`
		URL     string            `json:"url"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.URL == "" {
		http.Error(w, "name and url are required", http.StatusBadRequest)
		return
	}
	h.reg.Register(notify.WebhookTemplate{Name: body.Name, URL: body.URL, Headers: body.Headers})
	writeJSON(w, map[string]string{"status": "registered", "name": body.Name})
}

// handleTemplateRemove removes a named webhook template.
// DELETE /webhooks/templates?name=...
func (h *templateHandlers) handleTemplateRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if !h.reg.Remove(name) {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{"status": "removed", "name": name})
}

// handleTemplateList lists all registered webhook templates.
// GET /webhooks/templates
func (h *templateHandlers) handleTemplateList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]interface{}{"templates": h.reg.All()})
}
