package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExportStore(t *testing.T) *Store {
	t.Helper()
	s := NewStore(10)
	s.Record("nginx", Record{
		Time:       time.Now(),
		Kind:       "crash",
		Detail:     "exit code 1",
		CPUPercent: 0,
		MemMB:      0,
	})
	s.Record("nginx", Record{
		Time:       time.Now(),
		Kind:       "threshold",
		Detail:     "cpu high",
		CPUPercent: 95.2,
		MemMB:      120.5,
	})
	return s
}

func TestExportJSON_ValidOutput(t *testing.T) {
	s := makeExportStore(t)
	var buf bytes.Buffer
	if err := ExportJSON(s, &buf); err != nil {
		t.Fatalf("ExportJSON returned error: %v", err)
	}
	var result map[string][]Record
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result["nginx"]) != 2 {
		t.Errorf("expected 2 nginx records, got %d", len(result["nginx"]))
	}
}

func TestExportJSON_EmptyStore(t *testing.T) {
	s := NewStore(10)
	var buf bytes.Buffer
	if err := ExportJSON(s, &buf); err != nil {
		t.Fatalf("ExportJSON error on empty store: %v", err)
	}
	if !strings.Contains(buf.String(), "{}") {
		t.Errorf("expected empty JSON object, got: %s", buf.String())
	}
}

func TestExportText_ContainsHeaders(t *testing.T) {
	s := makeExportStore(t)
	var buf bytes.Buffer
	if err := ExportText(s, &buf); err != nil {
		t.Fatalf("ExportText returned error: %v", err)
	}
	out := buf.String()
	for _, header := range []string{"PROCESS", "TIME", "KIND", "CPU%", "MEM(MB)"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestExportText_ContainsData(t *testing.T) {
	s := makeExportStore(t)
	var buf bytes.Buffer
	_ = ExportText(s, &buf)
	out := buf.String()
	if !strings.Contains(out, "nginx") {
		t.Errorf("expected 'nginx' in text output")
	}
	if !strings.Contains(out, "crash") {
		t.Errorf("expected 'crash' in text output")
	}
}

func TestExportText_EmptyStore(t *testing.T) {
	s := NewStore(10)
	var buf bytes.Buffer
	if err := ExportText(s, &buf); err != nil {
		t.Fatalf("ExportText error on empty store: %v", err)
	}
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected no-records message")
	}
}
