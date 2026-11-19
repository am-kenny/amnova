package httpapi_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/am-kenny/amnova/internal/httpapi"
)

// okHandler answers every request with 200 and "ok"
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "ok")
})

// mustListen binds a free local port and fails the test if that errors
func mustListen(t *testing.T, addr string) *httpapi.Server {
	t.Helper()
	srv, err := httpapi.Listen(addr, okHandler)
	if err != nil {
		t.Fatalf("Listen(%s): %v", addr, err)
	}
	return srv
}

// start runs srv in the background and returns its cancel func and a channel with Run's result
func start(t *testing.T, srv *httpapi.Server) (context.CancelFunc, <-chan error) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx) }()
	t.Cleanup(cancel)
	return cancel, done
}

func TestServerServes(t *testing.T) {
	srv := mustListen(t, "127.0.0.1:0")
	start(t, srv)

	resp, err := http.Get("http://" + srv.Addr() + "/")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(body) != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
}

func TestServerShutsDownOnCancel(t *testing.T) {
	srv := mustListen(t, "127.0.0.1:0")
	cancel, done := start(t, srv)

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after cancel")
	}

	resp, err := http.Get("http://" + srv.Addr() + "/")
	if err == nil {
		_ = resp.Body.Close()
		t.Error("GET after shutdown: want an error, got nil")
	}
}

func TestListenTakenPort(t *testing.T) {
	first := mustListen(t, "127.0.0.1:0")

	if _, err := httpapi.Listen(first.Addr(), okHandler); err == nil {
		t.Error("Listen on a taken port: want an error, got nil")
	}
}
