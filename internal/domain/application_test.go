package domain_test

import (
	"errors"
	"maps"
	"testing"
	"time"

	"github.com/am-kenny/amnova/internal/domain"
)

// now is the fixed clock for these tests, 09:00 UTC given in Berlin time to check the UTC conversion
var now = time.Date(2026, 11, 25, 10, 0, 0, 0, time.FixedZone("CET", 60*60))

func TestValidateApplication(t *testing.T) {
	lastWeek := now.AddDate(0, 0, -7)

	tests := []struct {
		name          string
		in            domain.ApplicationInput
		wantAppliedAt string
	}{
		{"saved", domain.ApplicationInput{Status: domain.StatusSaved}, ""},
		{"applied now", domain.ApplicationInput{Status: domain.StatusApplied}, "2026-11-25T09:00:00Z"},
		{"applied last week", domain.ApplicationInput{Status: domain.StatusApplied, AppliedAt: &lastWeek}, "2026-11-18T09:00:00Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.Company = "  Aperture Science "
			tt.in.Role = "\tDevOps Engineer\n"

			got, err := tt.in.Validate(now)
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}

			if got.Company != "Aperture Science" || got.Role != "DevOps Engineer" {
				t.Errorf("company, role = %q, %q, want them trimmed", got.Company, got.Role)
			}

			var appliedAt string
			if got.AppliedAt != nil {
				appliedAt = got.AppliedAt.Format(time.RFC3339)
			}
			if appliedAt != tt.wantAppliedAt {
				t.Errorf("applied at = %q, want %q", appliedAt, tt.wantAppliedAt)
			}
		})
	}
}

func TestValidateApplicationErrors(t *testing.T) {
	tomorrow := now.AddDate(0, 0, 1)
	yesterday := now.AddDate(0, 0, -1)

	tests := []struct {
		name string
		in   domain.ApplicationInput
		want domain.FieldErrors
	}{
		{
			name: "everything missing",
			in:   domain.ApplicationInput{},
			want: domain.FieldErrors{"company": "required", "role": "required", "status": "must be saved or applied"},
		},
		{
			name: "blank names",
			in:   domain.ApplicationInput{Company: "  ", Role: "\t", Status: domain.StatusSaved},
			want: domain.FieldErrors{"company": "required", "role": "required"},
		},
		{
			name: "status a new application cannot start with",
			in:   domain.ApplicationInput{Company: "Aperture Science", Role: "DevOps Engineer", Status: domain.StatusOffer},
			want: domain.FieldErrors{"status": "must be saved or applied"},
		},
		{
			name: "applied in the future",
			in:   domain.ApplicationInput{Company: "Aperture Science", Role: "DevOps Engineer", Status: domain.StatusApplied, AppliedAt: &tomorrow},
			want: domain.FieldErrors{"applied_at": "cannot be in the future"},
		},
		{
			name: "saved with an applied date",
			in:   domain.ApplicationInput{Company: "Aperture Science", Role: "DevOps Engineer", Status: domain.StatusSaved, AppliedAt: &yesterday},
			want: domain.FieldErrors{"applied_at": "must be empty for a saved application"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.in.Validate(now)

			var got domain.FieldErrors
			if !errors.As(err, &got) {
				t.Fatalf("Validate error = %v, want FieldErrors", err)
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("Validate errors = %v, want %v", got, tt.want)
			}
		})
	}
}
