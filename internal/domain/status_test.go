package domain_test

import (
	"errors"
	"testing"

	"github.com/am-kenny/amnova/internal/domain"
)

func TestParseStatus(t *testing.T) {
	for _, want := range domain.Statuses {
		t.Run(string(want), func(t *testing.T) {
			got, err := domain.ParseStatus(string(want))
			if err != nil {
				t.Fatalf("ParseStatus(%q): %v", want, err)
			}
			if got != want {
				t.Errorf("ParseStatus(%q) = %q", want, got)
			}
		})
	}
}

func TestParseStatusInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"unknown", "wrong_status"},
		{"wrong case", "Applied"},
		{"label instead of value", "No response"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ParseStatus(tt.input)
			if !errors.Is(err, domain.ErrInvalidStatus) {
				t.Errorf("ParseStatus(%q) error = %v, want ErrInvalidStatus", tt.input, err)
			}
		})
	}
}
