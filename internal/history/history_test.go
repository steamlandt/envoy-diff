package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/envoy-diff/internal/diff"
	"github.com/example/envoy-diff/internal/history"
)

func makeEntry(src string, ts time.Time) history.Entry {
	return history.Entry{
		Timestamp: ts,
		SourceA:   src,
		SourceB:   "file:.env",
		Results: []diff.Entry{
			{Key: "FOO", ValA: "bar", ValB: "baz", Status: diff.Changed},
		},
	}
}

func TestRecord_And_List(t *testing.T) {
	dir := t.TempDir()

	e1 := makeEntry("file:a.env", time.Now().Add(-2*time.Hour))
	e2 := makeEntry("file:b.env", time.Now().Add(-1*time.Hour))

	if err := history.Record(dir, e1); err != nil {
		t.Fatalf("Record e1: %v", err)
	}
	if err := history.Record(dir, e2); err != nil {
		t.Fatalf("Record e2: %v", err)
	}

	entries, err := history.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Timestamp.Before(entries[1].Timestamp) {
		t.Error("expected entries sorted oldest-first")
	}
}

func TestList_MissingDir(t *testing.T) {
	entries, err := history.List("/nonexistent/path/xyz")
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}
}

func TestPrune_RemovesOldEntries(t *testing.T) {
	dir := t.TempDir()

	old := makeEntry("old", time.Now().Add(-48*time.Hour))
	recent := makeEntry("recent", time.Now().Add(-1*time.Hour))

	_ = history.Record(dir, old)
	_ = history.Record(dir, recent)

	if err := history.Prune(dir, 24*time.Hour); err != nil {
		t.Fatalf("Prune: %v", err)
	}

	entries, err := history.List(dir)
	if err != nil {
		t.Fatalf("List after prune: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after prune, got %d", len(entries))
	}
	if entries[0].SourceA != "recent" {
		t.Errorf("expected recent entry to survive, got %q", entries[0].SourceA)
	}
}

func TestRecord_CreatesParentDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/history"

	e := makeEntry("file:x.env", time.Now())
	if err := history.Record(dir, e); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(dir + "/history.json"); err != nil {
		t.Errorf("history.json not created: %v", err)
	}
}
