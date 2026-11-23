package health

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// readyTimeout bounds how long /readyz waits for its checks.
const readyTimeout = 2 * time.Second

// Pinger is a required /readyz check.
type Pinger interface {
	Ping(ctx context.Context) error
}

// NewHandler returns the probe routes
// /livez answers 200 while the process is able to serve requests.
// /readyz answers 200 while db answers its Ping, 503 otherwise.
func NewHandler(db Pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, _ *http.Request) {
		respond(w, http.StatusOK, "ok")
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), readyTimeout)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			slog.WarnContext(ctx, "readiness check failed", "err", err)
			respond(w, http.StatusServiceUnavailable, "not ready")
			return
		}

		respond(w, http.StatusOK, "ok")
	})

	return mux
}

func respond(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = fmt.Fprintln(w, body)
}
