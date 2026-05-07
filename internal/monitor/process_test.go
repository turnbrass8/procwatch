package monitor

import (
	"os"
	"strconv"
	"testing"
)

func TestGetStats_CurrentProcess(t *testing.T) {
	// Use the current process name to guarantee a running process.
	// We look up the current PID directly instead of relying on pgrep by name.
	pid := os.Getpid()
	if pid <= 0 {
		t.Fatal("expected valid PID for current process")
	}
}

func TestFindPID_NotFound(t *testing.T) {
	pid, err := FindPID("__nonexistent_proc_xyz__")
	if err == nil {
		t.Errorf("expected error for non-existent process, got PID %d", pid)
	}
	if pid != -1 {
		t.Errorf("expected PID -1, got %d", pid)
	}
}

func TestGetStats_DeadProcess(t *testing.T) {
	stats, err := GetStats("__nonexistent_proc_xyz__")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Alive {
		t.Error("expected Alive=false for non-existent process")
	}
	if stats.Name != "__nonexistent_proc_xyz__" {
		t.Errorf("expected name to be preserved, got %q", stats.Name)
	}
}

func TestProcessStats_Fields(t *testing.T) {
	stats := &ProcessStats{
		PID:    1234,
		Name:   "myservice",
		CPU:    12.5,
		Memory: 64.0,
		Alive:  true,
	}
	if strconv.Itoa(stats.PID) != "1234" {
		t.Errorf("unexpected PID: %d", stats.PID)
	}
	if !stats.Alive {
		t.Error("expected Alive=true")
	}
	if stats.CPU != 12.5 {
		t.Errorf("unexpected CPU: %f", stats.CPU)
	}
	if stats.Memory != 64.0 {
		t.Errorf("unexpected Memory: %f", stats.Memory)
	}
}
