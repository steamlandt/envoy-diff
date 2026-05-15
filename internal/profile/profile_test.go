package profile_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	return d
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := tempDir(t)
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := profile.Save(dir, "dev", env); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p, err := profile.Load(dir, "dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("Name = %q, want %q", p.Name, "dev")
	}
	if p.Env["FOO"] != "bar" || p.Env["BAZ"] != "qux" {
		t.Errorf("Env mismatch: %v", p.Env)
	}
	if p.SavedAt.IsZero() {
		t.Error("SavedAt should not be zero")
	}
	if time.Since(p.SavedAt) > 5*time.Second {
		t.Error("SavedAt seems too old")
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := tempDir(t)
	_, err := profile.Load(dir, "missing")
	if err != profile.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList_Empty(t *testing.T) {
	dir := tempDir(t)
	names, err := profile.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_MissingDir(t *testing.T) {
	names, err := profile.List("/nonexistent/path/xyz")
	if err != nil {
		t.Fatalf("List on missing dir should not error, got %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty, got %v", names)
	}
}

func TestList_Sorted(t *testing.T) {
	dir := tempDir(t)
	for _, n := range []string{"prod", "dev", "staging"} {
		if err := profile.Save(dir, n, map[string]string{"K": "v"}); err != nil {
			t.Fatal(err)
		}
	}
	names, err := profile.List(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"dev", "prod", "staging"}
	for i, w := range want {
		if names[i] != w {
			t.Errorf("names[%d] = %q, want %q", i, names[i], w)
		}
	}
}

func TestDelete_Existing(t *testing.T) {
	dir := tempDir(t)
	profile.Save(dir, "tmp", map[string]string{"X": "1"})
	if err := profile.Delete(dir, "tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tmp.json")); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestDelete_NotFound(t *testing.T) {
	dir := tempDir(t)
	if err := profile.Delete(dir, "ghost"); err != profile.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSave_CreatesParentDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "a", "b", "profiles")
	if err := profile.Save(dir, "x", map[string]string{}); err != nil {
		t.Fatalf("Save with deep path: %v", err)
	}
}
