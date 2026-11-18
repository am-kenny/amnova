package health

import (
	"fmt"
	"net/http"
)

// NewHandler returns the probe routes
// /livez answers 200 while the process is able to serve requests.
func NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, _ *http.Request) {
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
