package diff

import (
	"testing"
)

func makeStats(added, removed, changed, unchanged int) Stats {
	entries := make([]Entry, 0, added+removed+changed+unchanged)
	for i := 0; i < added; i++ {
		entries = append(entries, Entry{Status: StatusAdded})
	}
	for i := 0; i < removed; i++ {
		entries = append(entries, Entry{Status: StatusRemoved})
	}
	for i := 0; i < changed; i++ {
		entries = append(entries, Entry{Status: StatusChanged})
	}
	for i := 0; i < unchanged; i++ {
		entries = append(entries, Entry{Status: StatusUnchanged})
	}
	return Compute(entries)
}

func TestCompute_AllUnchanged(t *testing.T) {
	s := makeStats(0, 0, 0, 5)
	if s.Unchanged != 5 || s.Total != 5 || s.HasDiff() {
		t.Fatalf("unexpected stats: %+v", s)
	}
}

func TestCompute_Mixed(t *testing.T) {
	s := makeStats(2, 1, 3, 4)
	if s.Added != 2 || s.Removed != 1 || s.Changed != 3 || s.Unchanged != 4 {
		t.Fatalf("unexpected stats: %+v", s)
	}
	if s.Total != 10 {
		t.Fatalf("expected Total=10, got %d", s.Total)
	}
	if !s.HasDiff() {
		t.Fatal("expected HasDiff to be true")
	}
	if s.DiffCount() != 6 {
		t.Fatalf("expected DiffCount=6, got %d", s.DiffCount())
	}
}

func TestCompute_Empty(t *testing.T) {
	s := Compute(nil)
	if s.Total != 0 || s.HasDiff() {
		t.Fatalf("expected empty stats, got %+v", s)
	}
}

func TestFormatStats_NoDiff(t *testing.T) {
	s := makeStats(0, 0, 0, 3)
	got := FormatStats(s)
	if got != "no differences" {
		t.Fatalf("expected 'no differences', got %q", got)
	}
}

func TestFormatStats_OnlyAdded(t *testing.T) {
	s := makeStats(1, 0, 0, 0)
	got := FormatStats(s)
	if got != "1 added" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestFormatStats_Mixed(t *testing.T) {
	s := makeStats(2, 1, 3, 4)
	got := FormatStats(s)
	expected := "2 added, 1 removed, 3 changed, 4 unchanged"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestFormatStats_PluralSingular(t *testing.T) {
	s := makeStats(1, 1, 1, 1)
	got := FormatStats(s)
	expected := "1 added, 1 removed, 1 changed, 1 unchanged"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
