// Package summary provides aggregated statistics over a set of diff entries.
package summary

import (
	"fmt"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Stats holds counts of each change category from a diff result.
type Stats struct {
	Added   int
	Removed int
	Changed int
	Unchanged int
}

// Total returns the total number of keys across all categories.
func (s Stats) Total() int {
	return s.Added + s.Removed + s.Changed + s.Unchanged
}

// HasDiff reports whether any differences were found.
func (s Stats) HasDiff() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// Compute derives Stats from a slice of diff entries.
func Compute(entries []diff.Entry) Stats {
	var s Stats
	for _, e := range entries {
		switch e.Status {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Changed:
			s.Changed++
		case diff.Unchanged:
			s.Unchanged++
		}
	}
	return s
}

// Format returns a human-readable one-line summary string.
func Format(s Stats) string {
	if !s.HasDiff() {
		return fmt.Sprintf("No differences found (%d keys compared).", s.Total())
	}

	parts := make([]string, 0, 3)
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", s.Changed))
	}
	return fmt.Sprintf("%s (%d keys compared).", strings.Join(parts, ", "), s.Total())
}
