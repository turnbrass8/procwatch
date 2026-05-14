package notify

import (
	"strings"
	"sync"
)

// FilterRule defines a keyword-based suppression rule for a process.
type FilterRule struct {
	Process  string
	Keywords []string
}

// Filter suppresses alerts whose reason matches any registered keyword rule.
type Filter struct {
	mu    sync.RWMutex
	rules map[string][]string // process -> keywords
}

// NewFilter returns an empty Filter.
func NewFilter() *Filter {
	return &Filter{
		rules: make(map[string][]string),
	}
}

// AddRule registers keywords for a process name. Matching is case-insensitive.
func (f *Filter) AddRule(process string, keywords []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	existing := f.rules[process]
	f.rules[process] = append(existing, keywords...)
}

// RemoveRule deletes all keyword rules for a process.
func (f *Filter) RemoveRule(process string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.rules, process)
}

// IsSuppressed returns true if the alert reason matches any keyword for the process.
func (f *Filter) IsSuppressed(process, reason string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keywords, ok := f.rules[process]
	if !ok {
		return false
	}
	lower := strings.ToLower(reason)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// Rules returns a copy of all registered rules.
func (f *Filter) Rules() map[string][]string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	copy := make(map[string][]string, len(f.rules))
	for k, v := range f.rules {
		kws := make([]string, len(v))
		for i, kw := range v {
			kws[i] = kw
		}
		copy[k] = kws
	}
	return copy
}
