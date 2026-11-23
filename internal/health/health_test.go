package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/am-kenny/amnova/internal/health"
)

// pingerFunc turns a func into a health.Pinger
type pingerFunc func(context.Context) error

func (f pingerFunc) Ping(ctx context.Context) error { return f(ctx) }

var (
	healthy   = pingerFunc(func(context.Context) error { return nil })
	unhealthy = pingerFunc(func(context.Context) error { return errors.New("database is unavailable") })
)

// serve sends one request to the probe routes and returns the recorded response
func serve(db health.Pinger, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	health.NewHandler(db).ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

func TestLivez(t *testing.T) {
	// Liveness must not depend on the database
	rec := serve(unhealthy, http.MethodGet, "/livez")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "ok" {
		t.Errorf("body = %q, want %q", got, "ok")
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want %q", got, "no-store")
	}
}

func TestReadyz(t *testing.T) {
	tests := []struct {
		name       string
		db         health.Pinger
		wantStatus int
		wantBody   string
	}{
		{"ready", healthy, http.StatusOK, "ok"},
		{"database down", unhealthy, http.StatusServiceUnavailable, "not ready"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(tt.db, http.MethodGet, "/readyz")

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
			if got := rec.Header().Get("Cache-Control"); got != "no-store" {
				t.Errorf("Cache-Control = %q, want %q", got, "no-store")
			}
		})
	}
}

func TestReadyzTimesOut(t *testing.T) {
	// A ping that hangs must not hang the probe
	stuck := pingerFunc(func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	rec := serve(stuck, http.MethodGet, "/readyz")

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestRouting(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"livez", http.MethodGet, "/livez", http.StatusOK},
		{"livez rejects POST", http.MethodPost, "/livez", http.StatusMethodNotAllowed},
		{"readyz", http.MethodGet, "/readyz", http.StatusOK},
		{"readyz rejects POST", http.MethodPost, "/readyz", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(healthy, tt.method, tt.path)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
