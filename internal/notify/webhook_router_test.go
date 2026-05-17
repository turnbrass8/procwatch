package notify

import (
	"testing"
)

func newRegistry() *TemplateRegistry {
	return NewTemplateRegistry()
}

func TestTemplateRegistry_RegisterAndGet(t *testing.T) {
	reg := newRegistry()
	tmpl := WebhookTemplate{Name: "slack", URL: "https://hooks.slack.com/test", Headers: map[string]string{"X-Token": "abc"}}
	reg.Register(tmpl)

	got, ok := reg.Get("slack")
	if !ok {
		t.Fatal("expected template to be found")
	}
	if got.URL != tmpl.URL {
		t.Errorf("URL mismatch: got %q want %q", got.URL, tmpl.URL)
	}
	if got.Headers["X-Token"] != "abc" {
		t.Errorf("header mismatch")
	}
}

func TestTemplateRegistry_GetUnknown(t *testing.T) {
	reg := newRegistry()
	_, ok := reg.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown template")
	}
}

func TestTemplateRegistry_Overwrite(t *testing.T) {
	reg := newRegistry()
	reg.Register(WebhookTemplate{Name: "pagerduty", URL: "https://old.url"})
	reg.Register(WebhookTemplate{Name: "pagerduty", URL: "https://new.url"})

	got, _ := reg.Get("pagerduty")
	if got.URL != "https://new.url" {
		t.Errorf("expected updated URL, got %q", got.URL)
	}
}

func TestTemplateRegistry_Remove(t *testing.T) {
	reg := newRegistry()
	reg.Register(WebhookTemplate{Name: "teams", URL: "https://teams.example.com"})

	removed := reg.Remove("teams")
	if !removed {
		t.Error("expected Remove to return true")
	}
	_, ok := reg.Get("teams")
	if ok {
		t.Error("expected template to be gone after Remove")
	}
}

func TestTemplateRegistry_RemoveNonExistent(t *testing.T) {
	reg := newRegistry()
	if reg.Remove("ghost") {
		t.Error("expected Remove to return false for unknown name")
	}
}

func TestTemplateRegistry_All(t *testing.T) {
	reg := newRegistry()
	reg.Register(WebhookTemplate{Name: "a", URL: "https://a.example.com"})
	reg.Register(WebhookTemplate{Name: "b", URL: "https://b.example.com"})

	all := reg.All()
	if len(all) != 2 {
		t.Errorf("expected 2 templates, got %d", len(all))
	}
}

func TestTemplateRegistry_AllEmpty(t *testing.T) {
	reg := newRegistry()
	if len(reg.All()) != 0 {
		t.Error("expected empty slice for empty registry")
	}
}
