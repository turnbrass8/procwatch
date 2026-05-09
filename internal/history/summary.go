package history

import "time"

// ProcessSummary holds aggregated statistics for a single process.
type ProcessSummary struct {
	ProcessName    string
	TotalEvents    int
	CrashCount     int
	ThresholdCount int
	LastEvent      time.Time
	AvgCPU         float64
	AvgMemMB       float64
}

// Summarize computes a ProcessSummary from a slice of records.
func Summarize(name string, records []Record) ProcessSummary {
	if len(records) == 0 {
		return ProcessSummary{ProcessName: name}
	}

	summary := ProcessSummary{
		ProcessName: name,
		TotalEvents: len(records),
	}

	var totalCPU float64
	var totalMem float64

	for _, r := range records {
		switch r.Kind {
		case KindCrash:
			summary.CrashCount++
		case KindThreshold:
			summary.ThresholdCount++
		}
		if r.Timestamp.After(summary.LastEvent) {
			summary.LastEvent = r.Timestamp
		}
		totalCPU += r.CPUPercent
		totalMem += r.MemMB
	}

	n := float64(len(records))
	summary.AvgCPU = totalCPU / n
	summary.AvgMemMB = totalMem / n

	return summary
}

// SummarizeAll returns a summary for every process tracked in the store.
func SummarizeAll(s *Store) []ProcessSummary {
	all := s.All()
	summaries := make([]ProcessSummary, 0, len(all))
	for name, records := range all {
		summaries = append(summaries, Summarize(name, records))
	}
	return summaries
}
