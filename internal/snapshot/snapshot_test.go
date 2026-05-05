package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envoy-diff/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}

	if err := snapshot.Save(path, "test-label", "file://.env", env); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if snap.Label != "test-label" {
		t.Errorf("Label = %q, want %q", snap.Label, "test-label")
	}
	if snap.Source != "file://.env" {
		t.Errorf("Source = %q, want %q", snap.Source, "file://.env")
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
	if snap.CapturedAt.After(time.Now().Add(time.Second)) {
		t.Error("CapturedAt should not be in the future")
	}
	for k, v := range env {
		if snap.Env[k] != v {
			t.Errorf("Env[%q] = %q, want %q", k, snap.Env[k], v)
		}
	}
}

func TestSave_CreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "nested", "snap.json")

	if err := snapshot.Save(path, "lbl", "self", map[string]string{"A": "1"}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
