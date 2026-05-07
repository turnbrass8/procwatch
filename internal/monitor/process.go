package monitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ProcessStats holds resource usage for a single process.
type ProcessStats struct {
	PID    int
	Name   string
	CPU    float64 // percentage
	Memory float64 // MB
	Alive  bool
}

// FindPID returns the PID of a process by name, or -1 if not found.
func FindPID(name string) (int, error) {
	out, err := exec.Command("pgrep", "-x", name).Output()
	if err != nil {
		return -1, fmt.Errorf("process %q not found: %w", name, err)
	}
	pidStr := strings.TrimSpace(string(out))
	lines := strings.Split(pidStr, "\n")
	pid, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return -1, fmt.Errorf("failed to parse PID for %q: %w", name, err)
	}
	return pid, nil
}

// GetStats returns CPU and memory usage for a process by name.
func GetStats(name string) (*ProcessStats, error) {
	pid, err := FindPID(name)
	if err != nil {
		return &ProcessStats{Name: name, Alive: false}, nil
	}

	out, err := exec.Command(
		"ps", "-p", strconv.Itoa(pid), "-o", "pid,%cpu,rss", "--no-headers",
	).Output()
	if err != nil {
		return &ProcessStats{Name: name, PID: pid, Alive: false}, nil
	}

	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) < 3 {
		return nil, fmt.Errorf("unexpected ps output for %q", name)
	}

	cpu, _ := strconv.ParseFloat(fields[1], 64)
	rssKB, _ := strconv.ParseFloat(fields[2], 64)

	return &ProcessStats{
		PID:    pid,
		Name:   name,
		CPU:    cpu,
		Memory: rssKB / 1024.0,
		Alive:  true,
	}, nil
}
