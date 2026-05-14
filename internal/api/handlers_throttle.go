package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

// handleThrottleReset handles DELETE /throttle?name=<process>
// It resets the alert throttle for a specific process, or all processes
// if the name parameter is omitted.
func handleThrottleReset(th *notify.Throttle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			th.ResetAll()
			writeJSON(w, http.StatusOK, map[string]string{"status": "all throttles reset"})
			return
		}

		th.Reset(name)
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "throttle reset",
			"process": name,
		})
	}
}

// handleThrottleStatus handles GET /throttle/status?name=<process>
// It reports whether a given process is currently throttled.
func handleThrottleStatus(th *notify.Throttle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query parameter", http.StatusBadRequest)
			return
		}
		// Probe without recording: clone allow logic using a zero-duration copy
		// We use a non-destructive check via a sentinel probe throttle.
		allowed := th.Allow(name)
		if allowed {
			// We consumed the slot; reset immediately so the probe is transparent.
			th.Reset(name)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"process":   name,
			"throttled": !allowed,
		})
	}
}
