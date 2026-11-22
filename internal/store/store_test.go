package store

import (
	"os"
	"path/filepath"
	"testing"
)

// mustOpen opens a store in a fresh temporary dir and closes it when the test ends
func mustOpen(t *testing.T) (*Store, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "data")

	s, err := Open(t.Context(), dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	return s, dir
}

func TestOpenCreatesDatabase(t *testing.T) {
	_, dir := mustOpen(t)

	if _, err := os.Stat(filepath.Join(dir, fileName)); err != nil {
		t.Errorf("database file: %v", err)
	}
}

func TestOpenSetsPragmas(t *testing.T) {
	s, _ := mustOpen(t)

	tests := []struct {
		pragma string
		want   string
	}{
		{"journal_mode", "wal"},
		{"foreign_keys", "1"},
		{"busy_timeout", "5000"},
	}

	for _, tt := range tests {
		t.Run(tt.pragma, func(t *testing.T) {
			var got string
			if err := s.db.QueryRowContext(t.Context(), "PRAGMA "+tt.pragma).Scan(&got); err != nil {
				t.Fatalf("PRAGMA %s: %v", tt.pragma, err)
			}
			if got != tt.want {
				t.Errorf("%s = %q, want %q", tt.pragma, got, tt.want)
			}
		})
	}
}

func TestPing(t *testing.T) {
	s, _ := mustOpen(t)

	if err := s.Ping(t.Context()); err != nil {
		t.Errorf("Ping: %v", err)
	}
}
