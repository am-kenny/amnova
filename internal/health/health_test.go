package health_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/am-kenny/amnova/internal/health"
)

func TestLivez(t *testing.T) {
	rec := httptest.NewRecorder()
	health.NewHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/livez", nil))

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

func TestRouting(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"livez", http.MethodGet, "/livez", http.StatusOK},
		{"livez rejects POST", http.MethodPost, "/livez", http.StatusMethodNotAllowed},
		{"unknown path", http.MethodGet, "/", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			health.NewHandler().ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
