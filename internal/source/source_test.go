package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/source"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_File(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := source.Load(source.Source{Kind: source.File, Path: p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", env["FOO"], "bar")
	}
	if env["BAZ"] != "qux" {
		t.Errorf("BAZ: got %q, want %q", env["BAZ"], "qux")
	}
}

func TestLoad_File_EmptyPath(t *testing.T) {
	_, err := source.Load(source.Source{Kind: source.File})
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestLoad_Self(t *testing.T) {
	t.Setenv("ENVOY_DIFF_TEST_KEY", "hello")
	env, err := source.Load(source.Source{Kind: source.Self})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["ENVOY_DIFF_TEST_KEY"] != "hello" {
		t.Errorf("ENVOY_DIFF_TEST_KEY: got %q, want %q", env["ENVOY_DIFF_TEST_KEY"], "hello")
	}
}

func TestLoad_PID_Invalid(t *testing.T) {
	_, err := source.Load(source.Source{Kind: source.PID, PID: -1})
	if err == nil {
		t.Fatal("expected error for invalid PID, got nil")
	}
}

func TestLoad_UnknownKind(t *testing.T) {
	_, err := source.Load(source.Source{Kind: "s3"})
	if err == nil {
		t.Fatal("expected error for unknown kind, got nil")
	}
}

func TestSource_String(t *testing.T) {
	cases := []struct {
		src  source.Source
		want string
	}{
		{source.Source{Kind: source.File, Path: ".env"}, "file:.env"},
		{source.Source{Kind: source.PID, PID: 42}, "pid:42"},
		{source.Source{Kind: source.Self}, "self"},
		{source.Source{Kind: "other"}, "unknown(other)"},
	}
	for _, tc := range cases {
		if got := tc.src.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
