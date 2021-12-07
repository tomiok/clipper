package clipper

import (
	"encoding/json"
	"net/http"
)

// ExposeMetrics could be used to see the metrics given the command name.
// This function should be used as a handler or trying to fit in other context in web frameworks.
func ExposeMetrics() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		cmd := values.Get("command")

		if cmd == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		s := FillStats(cmd, false)

		_ = json.NewEncoder(w).Encode(&s)
	}
}
