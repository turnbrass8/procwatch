package notify

import (
	"sync"
)

// WebhookTemplate holds a named webhook endpoint with optional metadata.
type WebhookTemplate struct {
	Name    string
	URL     string
	Headers map[string]string
}

// TemplateRegistry stores named webhook templates that can be referenced
// by routing rules instead of hardcoding URLs.
type TemplateRegistry struct {
	mu        sync.RWMutex
	templates map[string]WebhookTemplate
}

// NewTemplateRegistry creates an empty TemplateRegistry.
func NewTemplateRegistry() *TemplateRegistry {
	return &TemplateRegistry{
		templates: make(map[string]WebhookTemplate),
	}
}

// Register adds or replaces a named webhook template.
func (r *TemplateRegistry) Register(t WebhookTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[t.Name] = t
}

// Remove deletes a named webhook template. Returns true if it existed.
func (r *TemplateRegistry) Remove(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.templates[name]
	if ok {
		delete(r.templates, name)
	}
	return ok
}

// Get retrieves a template by name. Returns false if not found.
func (r *TemplateRegistry) Get(name string) (WebhookTemplate, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.templates[name]
	return t, ok
}

// All returns a snapshot of all registered templates.
func (r *TemplateRegistry) All() []WebhookTemplate {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]WebhookTemplate, 0, len(r.templates))
	for _, t := range r.templates {
		out = append(out, t)
	}
	return out
}
