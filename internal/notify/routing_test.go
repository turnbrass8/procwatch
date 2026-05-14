package notify

import (
	"testing"
)

func newRouter() *Router {
	return NewRouter("http://default.example.com/webhook")
}

func TestRouter_DefaultURLWhenNoRoutes(t *testing.T) {
	r := newRouter()
	urls := r.Resolve("nginx")
	if len(urls) != 1 || urls[0] != "http://default.example.com/webhook" {
		t.Fatalf("expected default URL, got %v", urls)
	}
}

func TestRouter_MatchesSpecificProcess(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{
		Name:       "nginx-route",
		WebhookURL: "http://nginx.example.com/hook",
		Processes:  []string{"nginx"},
	})

	urls := r.Resolve("nginx")
	if len(urls) != 1 || urls[0] != "http://nginx.example.com/hook" {
		t.Fatalf("expected nginx-specific URL, got %v", urls)
	}
}

func TestRouter_FallsBackForUnmatchedProcess(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{
		Name:       "nginx-route",
		WebhookURL: "http://nginx.example.com/hook",
		Processes:  []string{"nginx"},
	})

	urls := r.Resolve("redis")
	if len(urls) != 1 || urls[0] != "http://default.example.com/webhook" {
		t.Fatalf("expected default URL for unmatched process, got %v", urls)
	}
}

func TestRouter_CatchAllRouteMatchesAnyProcess(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{
		Name:       "catch-all",
		WebhookURL: "http://all.example.com/hook",
		Processes:  []string{},
	})

	urls := r.Resolve("anything")
	if len(urls) != 1 || urls[0] != "http://all.example.com/hook" {
		t.Fatalf("expected catch-all URL, got %v", urls)
	}
}

func TestRouter_MultipleMatchingRoutes(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{Name: "r1", WebhookURL: "http://a.example.com", Processes: []string{"nginx"}})
	r.AddRoute(Route{Name: "r2", WebhookURL: "http://b.example.com", Processes: []string{"nginx", "redis"}})

	urls := r.Resolve("nginx")
	if len(urls) != 2 {
		t.Fatalf("expected 2 matching URLs, got %d: %v", len(urls), urls)
	}
}

func TestRouter_RemoveRoute(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{Name: "nginx-route", WebhookURL: "http://nginx.example.com/hook", Processes: []string{"nginx"}})

	removed := r.RemoveRoute("nginx-route")
	if !removed {
		t.Fatal("expected route to be removed")
	}

	urls := r.Resolve("nginx")
	if len(urls) != 1 || urls[0] != "http://default.example.com/webhook" {
		t.Fatalf("expected fallback to default after removal, got %v", urls)
	}
}

func TestRouter_RemoveRoute_NotFound(t *testing.T) {
	r := newRouter()
	removed := r.RemoveRoute("nonexistent")
	if removed {
		t.Fatal("expected false for nonexistent route")
	}
}

func TestRouter_Routes_ReturnsCopy(t *testing.T) {
	r := newRouter()
	r.AddRoute(Route{Name: "r1", WebhookURL: "http://a.example.com"})
	routes := r.Routes()
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	// Mutating the returned slice should not affect internal state
	routes[0].Name = "mutated"
	if r.Routes()[0].Name != "r1" {
		t.Fatal("Routes() should return a copy, not a reference")
	}
}
