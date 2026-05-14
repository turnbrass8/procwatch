package api

import (
	"net/http"

	"github.com/user/procwatch/internal/notify"
)

func handleDigestFlush(d *notify.Digest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query param", http.StatusBadRequest)
			return
		}

		d.Flush(name)
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "flushed",
			"process": name,
		})
	}
}

func handleDigestStatus(d *notify.Digest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing 'name' query param", http.StatusBadRequest)
			return
		}

		pending := d.PendingCount(name)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"process": name,
			"pending": pending,
		})
	}
}
