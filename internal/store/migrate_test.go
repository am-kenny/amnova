package store

import "testing"

const ts = "2026-11-23T10:00:00Z"

// mustMigrate opens a store in a fresh temporary dir and applies all migrations
func mustMigrate(t *testing.T) *Store {
	t.Helper()
	s, _ := mustOpen(t)

	if err := s.Migrate(t.Context()); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	return s
}

func TestMigrateTwice(t *testing.T) {
	s := mustMigrate(t)

	if err := s.Migrate(t.Context()); err != nil {
		t.Errorf("second Migrate: %v", err)
	}
}

func TestApplicationConstraints(t *testing.T) {
	s := mustMigrate(t)

	_, err := s.db.ExecContext(t.Context(),
		"INSERT INTO companies (id, name, created_at, updated_at) VALUES (1, 'Acme', ?, ?)", ts, ts,
	)
	if err != nil {
		t.Fatalf("insert company: %v", err)
	}

	// A valid application with every field set
	_, err = s.db.ExecContext(t.Context(), `
		INSERT INTO applications (id, company_id, role, status, applied_at, created_at, updated_at)
		VALUES (1, 1, 'DevOps Engineer', 'applied', ?, ?, ?)`, ts, ts, ts,
	)
	if err != nil {
		t.Fatalf("insert valid application: %v", err)
	}

	// Each change below breaks one rule and must be rejected
	tests := []struct {
		name string
		set  string
	}{
		{"unknown status", "status = 'bad_status'"},
		{"missing company", "company_id = 99"},
		{"missing role", "role = NULL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.db.ExecContext(t.Context(), "UPDATE applications SET "+tt.set+" WHERE id = 1")
			if err == nil {
				t.Errorf("SET %s: want an error, got nil", tt.set)
			}
		})
	}
}

func TestCompanyNameIsUniqueIgnoringCase(t *testing.T) {
	s := mustMigrate(t)

	const insert = "INSERT INTO companies (name, created_at, updated_at) VALUES (?, ?, ?)"

	if _, err := s.db.ExecContext(t.Context(), insert, "Acme", ts, ts); err != nil {
		t.Fatalf("insert Acme: %v", err)
	}
	if _, err := s.db.ExecContext(t.Context(), insert, "ACME", ts, ts); err == nil {
		t.Error("insert ACME: want an error, got nil")
	}
}
