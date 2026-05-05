package summary_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/summary"
)

func makeEntries(statuses ...diff.Status) []diff.Entry {
	entries := make([]diff.Entry, len(statuses))
	for i, s := range statuses {
		entries[i] = diff.Entry{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return entries
}

func TestCompute_AllUnchanged(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", Status: diff.Unchanged},
		{Key: "B", Status: diff.Unchanged},
	}
	s := summary.Compute(entries)
	if s.Unchanged != 2 || s.Added != 0 || s.Removed != 0 || s.Changed != 0 {
		t.Errorf("unexpected stats: %+v", s)
	}
	if s.HasDiff() {
		t.Error("expected HasDiff to be false")
	}
	if s.Total() != 2 {
		t.Errorf("expected total 2, got %d", s.Total())
	}
}

func TestCompute_Mixed(t *testing.T) {
	entries := []diff.Entry{
		{Key: "A", Status: diff.Added},
		{Key: "B", Status: diff.Removed},
		{Key: "C", Status: diff.Changed},
		{Key: "D", Status: diff.Unchanged},
	}
	s := summary.Compute(entries)
	if s.Added != 1 || s.Removed != 1 || s.Changed != 1 || s.Unchanged != 1 {
		t.Errorf("unexpected stats: %+v", s)
	}
	if !s.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}

func TestFormat_NoDiff(t *testing.T) {
	s := summary.Stats{Unchanged: 5}
	out := summary.Format(s)
	expected := "No differences found (5 keys compared)."
	if out != expected {
		t.Errorf("got %q, want %q", out, expected)
	}
}

func TestFormat_WithDiff(t *testing.T) {
	s := summary.Stats{Added: 2, Removed: 1, Changed: 3, Unchanged: 4}
	out := summary.Format(s)
	if out == "" {
		t.Fatal("expected non-empty format output")
	}
	for _, sub := range []string{"2 added", "1 removed", "3 changed", "10 keys compared"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected %q in output %q", sub, out)
		}
	}
}

func TestFormat_OnlyAdded(t *testing.T) {
	s := summary.Stats{Added: 3}
	out := summary.Format(s)
	if !strings.Contains(out, "3 added") {
		t.Errorf("unexpected output: %q", out)
	}
	if strings.Contains(out, "removed") || strings.Contains(out, "changed") {
		t.Errorf("should not mention removed/changed: %q", out)
	}
}
