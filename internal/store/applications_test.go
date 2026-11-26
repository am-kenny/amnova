package store

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/am-kenny/amnova/internal/domain"
)

// created is the clock for these tests, given in Berlin time to check the UTC conversion
var created = time.Date(2026, 11, 25, 10, 0, 0, 0, time.FixedZone("CET", 60*60))

func TestCreateApplication(t *testing.T) {
	s := mustMigrate(t)
	appliedAt := created.Add(-time.Hour)

	id, err := s.CreateApplication(t.Context(), domain.ApplicationInput{
		Company:   "Aperture Science",
		Role:      "DevOps Engineer",
		Status:    domain.StatusApplied,
		AppliedAt: &appliedAt,
	}, created)
	if err != nil {
		t.Fatalf("CreateApplication: %v", err)
	}

	// Times are stored as RFC 3339 in UTC, so they sort correctly as text
	var company, role, status, applied, createdAt string
	err = s.db.QueryRowContext(t.Context(), `
		SELECT c.name, a.role, a.status, a.applied_at, a.created_at
		FROM applications a JOIN companies c ON c.id = a.company_id
		WHERE a.id = ?`, id,
	).Scan(&company, &role, &status, &applied, &createdAt)
	if err != nil {
		t.Fatalf("read application: %v", err)
	}

	got := []string{company, role, status, applied, createdAt}
	want := []string{"Aperture Science", "DevOps Engineer", "applied", "2026-11-25T08:00:00Z", "2026-11-25T09:00:00Z"}
	if !slices.Equal(got, want) {
		t.Errorf("stored %v, want %v", got, want)
	}
}

func TestCreateAndListApplications(t *testing.T) {
	s := mustMigrate(t)
	later := created.Add(time.Hour).UTC()

	inputs := []domain.ApplicationInput{
		{Company: "Aperture Science", Role: "DevOps Engineer", Status: domain.StatusSaved},
		{Company: "Black Mesa", Role: "Platform Engineer", Status: domain.StatusApplied, AppliedAt: &later},
	}
	clocks := []time.Time{created, later}

	ids := make([]int64, len(inputs))
	for i, in := range inputs {
		id, err := s.CreateApplication(t.Context(), in, clocks[i])
		if err != nil {
			t.Fatalf("CreateApplication(%s): %v", in.Company, err)
		}
		ids[i] = id
	}

	got, err := s.ListApplications(t.Context())
	if err != nil {
		t.Fatalf("ListApplications: %v", err)
	}

	// Newest first, times back in UTC
	want := []domain.Application{
		{ID: ids[1], Company: "Black Mesa", Role: "Platform Engineer", Status: domain.StatusApplied, AppliedAt: &later, CreatedAt: later},
		{ID: ids[0], Company: "Aperture Science", Role: "DevOps Engineer", Status: domain.StatusSaved, CreatedAt: created.UTC()},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListApplications =\n%+v\nwant\n%+v", got, want)
	}
}

func TestCreateApplicationIsAtomic(t *testing.T) {
	s := mustMigrate(t)

	// An unknown status passes the company upsert but fails the insert
	_, err := s.CreateApplication(t.Context(), domain.ApplicationInput{
		Company: "Black Mesa",
		Role:    "DevOps Engineer",
		Status:  "unknown_status",
	}, created)
	if err == nil {
		t.Fatal("CreateApplication: want an error, got nil")
	}

	var n int
	if err := s.db.QueryRowContext(t.Context(), "SELECT count(*) FROM companies").Scan(&n); err != nil {
		t.Fatalf("count companies: %v", err)
	}
	if n != 0 {
		t.Errorf("companies = %d, want 0: the failed insert must roll back the company", n)
	}
}
