package config_test

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/am-kenny/amnova/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantAddr  string
		wantLevel slog.Level
	}{
		{"defaults", nil, ":8080", slog.LevelInfo},
		{"listen addr", map[string]string{"AMNOVA_LISTEN_ADDR": "127.0.0.1:9000"}, "127.0.0.1:9000", slog.LevelInfo},
		{"debug level", map[string]string{"AMNOVA_LOG_LEVEL": "debug"}, ":8080", slog.LevelDebug},
		{"level is case insensitive", map[string]string{"AMNOVA_LOG_LEVEL": "WARN"}, ":8080", slog.LevelWarn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("Load: %v", err)
			}

			if cfg.Server.ListenAddr != tt.wantAddr {
				t.Errorf("ListenAddr = %q, want %q", cfg.Server.ListenAddr, tt.wantAddr)
			}
			if cfg.Log.Level != tt.wantLevel {
				t.Errorf("Level = %v, want %v", cfg.Log.Level, tt.wantLevel)
			}
		})
	}
}

func TestLoadInvalidLevel(t *testing.T) {
	clearEnv(t)
	t.Setenv("AMNOVA_LOG_LEVEL", "loud")

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load: want an error, got nil")
	}
	if !strings.Contains(err.Error(), "AMNOVA_LOG_LEVEL") {
		t.Errorf("error = %q, want it to name AMNOVA_LOG_LEVEL", err)
	}
}

// clearEnv unsets every variable Load reads, so the host environment cannot leak into a test
func clearEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{"AMNOVA_LISTEN_ADDR", "AMNOVA_LOG_LEVEL"} {
		t.Setenv(name, "")
	}
}
