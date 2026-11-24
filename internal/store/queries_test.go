package store

import (
	"testing"

	"github.com/am-kenny/amnova/internal/store/sqlc"
)

func TestUpsertCompanyMatchesNameIgnoringCase(t *testing.T) {
	s := mustMigrate(t)
	q := sqlc.New(s.db)

	upsert := func(name string) int64 {
		t.Helper()
		id, err := q.UpsertCompany(t.Context(), sqlc.UpsertCompanyParams{Name: name, CreatedAt: ts, UpdatedAt: ts})
		if err != nil {
			t.Fatalf("upsert %s: %v", name, err)
		}
		return id
	}

	first := upsert("Aperture Science")
	if second := upsert("APERTURE SCIENCE"); second != first {
		t.Errorf("APERTURE SCIENCE id = %d, want %d (the existing Aperture Science)", second, first)
	}

	var name string
	if err := s.db.QueryRowContext(t.Context(), "SELECT name FROM companies WHERE id = ?", first).Scan(&name); err != nil {
		t.Fatalf("read name: %v", err)
	}
	if name != "Aperture Science" {
		t.Errorf("name = %q, want the first spelling %q", name, "Aperture Science")
	}
}
