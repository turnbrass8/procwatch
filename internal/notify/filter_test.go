package notify

import (
	"testing"
)

func newFilter() *Filter { return NewFilter() }

func TestFilter_AllowsWithNoRules(t *testing.T) {
	f := newFilter()
	if f.IsSuppressed("nginx", "crashed unexpectedly") {
		t.Fatal("expected alert to be allowed with no rules")
	}
}

func TestFilter_SuppressesMatchingKeyword(t *testing.T) {
	f := newFilter()
	f.AddRule("nginx", []string{"oom", "timeout"})
	if !f.IsSuppressed("nginx", "process hit OOM killer") {
		t.Fatal("expected alert to be suppressed by 'oom' keyword")
	}
}

func TestFilter_CaseInsensitiveMatch(t *testing.T) {
	f := newFilter()
	f.AddRule("redis", []string{"Timeout"})
	if !f.IsSuppressed("redis", "connection timeout exceeded") {
		t.Fatal("expected case-insensitive match to suppress alert")
	}
}

func TestFilter_AllowsNonMatchingReason(t *testing.T) {
	f := newFilter()
	f.AddRule("redis", []string{"oom"})
	if f.IsSuppressed("redis", "disk full") {
		t.Fatal("expected non-matching reason to be allowed")
	}
}

func TestFilter_IndependentPerProcess(t *testing.T) {
	f := newFilter()
	f.AddRule("nginx", []string{"oom"})
	if f.IsSuppressed("redis", "oom") {
		t.Fatal("expected rule for nginx not to affect redis")
	}
}

func TestFilter_RemoveRule(t *testing.T) {
	f := newFilter()
	f.AddRule("nginx", []string{"crash"})
	f.RemoveRule("nginx")
	if f.IsSuppressed("nginx", "crash detected") {
		t.Fatal("expected alert to be allowed after rule removal")
	}
}

func TestFilter_Rules_ReturnsCopy(t *testing.T) {
	f := newFilter()
	f.AddRule("svc", []string{"panic"})
	rules := f.Rules()
	rules["svc"] = append(rules["svc"], "injected")
	if f.IsSuppressed("svc", "injected") {
		t.Fatal("expected Rules() to return a copy, not a reference")
	}
}
