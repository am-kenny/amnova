package domain

import (
	"maps"
	"slices"
	"strings"
)

// FieldErrors maps an input field to what is wrong with it.
type FieldErrors map[string]string

func (e FieldErrors) Error() string {
	parts := make([]string, 0, len(e))
	for _, field := range slices.Sorted(maps.Keys(e)) {
		parts = append(parts, field+": "+e[field])
	}

	return strings.Join(parts, "; ")
}
