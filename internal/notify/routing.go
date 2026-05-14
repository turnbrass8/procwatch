package notify

import "sync"

// Route defines a named webhook destination with an optional process filter.
type Route struct {
	Name       string
	WebhookURL string
	// Processes is an allow-list of process names. Empty means all processes.
	Processes []string
}

// Router maps alert events to one or more webhook destinations based on
// configured routes. If no route matches, the default URL is used.
type Router struct {
	mu         sync.RWMutex
	routes     []Route
	defaultURL string
}

// NewRouter creates a Router with the given default webhook URL.
func NewRouter(defaultURL string) *Router {
	return &Router{defaultURL: defaultURL}
}

// AddRoute appends a route to the router.
func (r *Router) AddRoute(route Route) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
}

// RemoveRoute removes the first route whose Name matches the given name.
// Returns true if a route was removed.
func (r *Router) RemoveRoute(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rt := range r.routes {
		if rt.Name == name {
			r.routes = append(r.routes[:i], r.routes[i+1:]...)
			return true
		}
	}
	return false
}

// Resolve returns all webhook URLs that should receive an alert for the given
// process name. Falls back to the default URL when no specific route matches.
func (r *Router) Resolve(processName string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []string
	for _, rt := range r.routes {
		if len(rt.Processes) == 0 {
			matched = append(matched, rt.WebhookURL)
			continue
		}
		for _, p := range rt.Processes {
			if p == processName {
				matched = append(matched, rt.WebhookURL)
				break
			}
		}
	}

	if len(matched) == 0 {
		return []string{r.defaultURL}
	}
	return matched
}

// Routes returns a copy of the current route list.
func (r *Router) Routes() []Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Route, len(r.routes))
	copy(out, r.routes)
	return out
}
