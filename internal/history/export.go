package history

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// ExportJSON writes all history records as JSON to the given writer.
func ExportJSON(store *Store, w io.Writer) error {
	all := store.All()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(all); err != nil {
		return fmt.Errorf("history: json encode failed: %w", err)
	}
	return nil
}

// ExportText writes a human-readable summary table to the given writer.
func ExportText(store *Store, w io.Writer) error {
	all := store.All()
	if len(all) == 0 {
		_, err := fmt.Fprintln(w, "No history records found.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROCESS\tTIME\tKIND\tCPU%\tMEM(MB)\tDETAIL")
	fmt.Fprintln(tw, "-------\t----\t----\t----\t-------\t------")

	for process, records := range all {
		for _, r := range records {
			kind := "crash"
			if r.Kind == "threshold" {
				kind = "threshold"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%.1f\t%.1f\t%s\n",
				process,
				r.Time.Format(time.RFC3339),
				kind,
				r.CPUPercent,
				r.MemMB,
				r.Detail,
			)
		}
	}

	return tw.Flush()
}
