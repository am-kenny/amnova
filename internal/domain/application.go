package domain

import (
	"strings"
	"time"
)

// ApplicationInput is the data entered for a new application.
type ApplicationInput struct {
	Company   string
	Role      string
	Status    Status
	AppliedAt *time.Time
}

// Validate checks the input and returns a cleaned copy.
// It returns FieldErrors listing every problem.
func (in ApplicationInput) Validate(now time.Time) (ApplicationInput, error) {
	now = now.UTC()
	errs := FieldErrors{}

	out := ApplicationInput{
		Company: strings.TrimSpace(in.Company),
		Role:    strings.TrimSpace(in.Role),
		Status:  in.Status,
	}

	if out.Company == "" {
		errs["company"] = "required"
	}
	if out.Role == "" {
		errs["role"] = "required"
	}

	switch in.Status {
	case StatusSaved:
		if in.AppliedAt != nil {
			errs["applied_at"] = "must be empty for a saved application"
		}
	case StatusApplied:
		appliedAt := now
		if in.AppliedAt != nil {
			appliedAt = in.AppliedAt.UTC()
		}
		if appliedAt.After(now) {
			errs["applied_at"] = "cannot be in the future"
		}
		out.AppliedAt = &appliedAt
	default:
		errs["status"] = "must be saved or applied"
	}

	if len(errs) > 0 {
		return ApplicationInput{}, errs
	}

	return out, nil
}

// Application is a stored application.
type Application struct {
	ID        int64
	Company   string
	Role      string
	Status    Status
	AppliedAt *time.Time
	CreatedAt time.Time
}
