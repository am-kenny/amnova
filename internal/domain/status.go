package domain

import (
	"errors"
	"fmt"
	"slices"
)

// Status is where an application stands in the hiring process.
type Status string

const (
	StatusSaved        Status = "saved"
	StatusApplied      Status = "applied"
	StatusInterviewing Status = "interviewing"
	StatusOffer        Status = "offer"
	StatusAccepted     Status = "accepted"
	StatusRejected     Status = "rejected"
	StatusNoResponse   Status = "no_response"
)

// Statuses lists every status in board order.
var Statuses = []Status{
	StatusSaved,
	StatusApplied,
	StatusInterviewing,
	StatusOffer,
	StatusAccepted,
	StatusRejected,
	StatusNoResponse,
}

// ErrInvalidStatus is returned for a value that is not a known status.
var ErrInvalidStatus = errors.New("invalid status")

// ParseStatus returns the status named s, or ErrInvalidStatus.
func ParseStatus(s string) (Status, error) {
	st := Status(s)
	if !slices.Contains(Statuses, st) {
		return "", fmt.Errorf("%w: %q", ErrInvalidStatus, s)
	}

	return st, nil
}
