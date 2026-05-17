package notify

import (
	"strings"
	"sync"
)

// Tag represents a key-value label attached to a process alert context.
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Tagger manages per-process tags used to enrich alert payloads.
type Tagger struct {
	mu   sync.RWMutex
	tags map[string][]Tag // keyed by process name
}

// NewTagger creates a new Tagger instance.
func NewTagger() *Tagger {
	return &Tagger{
		tags: make(map[string][]Tag),
	}
}

// Set adds or replaces a tag for the given process name.
// If a tag with the same key already exists it is overwritten.
func (t *Tagger) Set(process, key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	existing := t.tags[process]
	for i, tag := range existing {
		if strings.EqualFold(tag.Key, key) {
			existing[i].Value = value
			t.tags[process] = existing
			return
		}
	}
	t.tags[process] = append(existing, Tag{Key: key, Value: value})
}

// Remove deletes a tag by key for the given process.
func (t *Tagger) Remove(process, key string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	existing := t.tags[process]
	updated := existing[:0]
	for _, tag := range existing {
		if !strings.EqualFold(tag.Key, key) {
			updated = append(updated, tag)
		}
	}
	t.tags[process] = updated
}

// Get returns all tags for the given process.
func (t *Tagger) Get(process string) []Tag {
	t.mu.RLock()
	defer t.mu.RUnlock()

	copy := make([]Tag, len(t.tags[process]))
	copy_ := t.tags[process]
	result := make([]Tag, len(copy_))
	for i, tag := range copy_ {
		result[i] = tag
	}
	return result
}

// Clear removes all tags for the given process, or all processes if empty string.
func (t *Tagger) Clear(process string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if process == "" {
		t.tags = make(map[string][]Tag)
		return
	}
	delete(t.tags, process)
}
