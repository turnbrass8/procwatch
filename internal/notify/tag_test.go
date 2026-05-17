package notify

import (
	"testing"
)

func newTagger() *Tagger {
	return NewTagger()
}

func TestTagger_SetAndGet(t *testing.T) {
	tg := newTagger()
	tg.Set("nginx", "env", "production")
	tg.Set("nginx", "team", "platform")

	tags := tg.Get("nginx")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestTagger_OverwritesSameKey(t *testing.T) {
	tg := newTagger()
	tg.Set("redis", "env", "staging")
	tg.Set("redis", "env", "production")

	tags := tg.Get("redis")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag after overwrite, got %d", len(tags))
	}
	if tags[0].Value != "production" {
		t.Errorf("expected value 'production', got %q", tags[0].Value)
	}
}

func TestTagger_CaseInsensitiveKey(t *testing.T) {
	tg := newTagger()
	tg.Set("app", "ENV", "dev")
	tg.Set("app", "env", "prod")

	tags := tg.Get("app")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag (case-insensitive dedup), got %d", len(tags))
	}
	if tags[0].Value != "prod" {
		t.Errorf("expected 'prod', got %q", tags[0].Value)
	}
}

func TestTagger_Remove(t *testing.T) {
	tg := newTagger()
	tg.Set("worker", "region", "us-east-1")
	tg.Set("worker", "tier", "free")
	tg.Remove("worker", "tier")

	tags := tg.Get("worker")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag after remove, got %d", len(tags))
	}
	if tags[0].Key != "region" {
		t.Errorf("expected remaining key 'region', got %q", tags[0].Key)
	}
}

func TestTagger_GetUnknownProcess(t *testing.T) {
	tg := newTagger()
	tags := tg.Get("unknown")
	if len(tags) != 0 {
		t.Errorf("expected empty slice for unknown process, got %d tags", len(tags))
	}
}

func TestTagger_ClearNamed(t *testing.T) {
	tg := newTagger()
	tg.Set("svc-a", "k", "v")
	tg.Set("svc-b", "k", "v")
	tg.Clear("svc-a")

	if len(tg.Get("svc-a")) != 0 {
		t.Error("expected svc-a tags to be cleared")
	}
	if len(tg.Get("svc-b")) != 1 {
		t.Error("expected svc-b tags to remain")
	}
}

func TestTagger_ClearAll(t *testing.T) {
	tg := newTagger()
	tg.Set("svc-a", "k", "v")
	tg.Set("svc-b", "k", "v")
	tg.Clear("")

	if len(tg.Get("svc-a")) != 0 || len(tg.Get("svc-b")) != 0 {
		t.Error("expected all tags cleared")
	}
}
