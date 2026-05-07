package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "procwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://hooks.example.com/abc
interval: 30s
processes:
  - name: nginx
    match: nginx
    max_cpu: 80
    max_mem_mb: 512
    alert_on_exit: true
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WebhookURL != "https://hooks.example.com/abc" {
		t.Errorf("webhook_url mismatch: %q", cfg.WebhookURL)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("interval mismatch: %v", cfg.Interval)
	}
	if len(cfg.Processes) != 1 {
		t.Fatalf("expected 1 process, got %d", len(cfg.Processes))
	}
	p := cfg.Processes[0]
	if p.Name != "nginx" || p.Match != "nginx" {
		t.Errorf("process fields mismatch: %+v", p)
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://hooks.example.com/abc
processes:
  - name: redis
    match: redis-server
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 15*time.Second {
		t.Errorf("expected default interval 15s, got %v", cfg.Interval)
	}
}

func TestLoad_MissingWebhook(t *testing.T) {
	path := writeTempConfig(t, `
processes:
  - name: app
    match: myapp
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing webhook_url, got nil")
	}
}

func TestLoad_MissingProcesses(t *testing.T) {
	path := writeTempConfig(t, `webhook_url: https://hooks.example.com/abc\n`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty processes, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
