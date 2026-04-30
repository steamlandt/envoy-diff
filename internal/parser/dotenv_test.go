package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := ParseFile(path)
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

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SINGLE='hello world'
DOUBLE="hello world"
`)
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SINGLE"] != "hello world" {
		t.Errorf("SINGLE: got %q", env["SINGLE"])
	}
	if env["DOUBLE"] != "hello world" {
		t.Errorf("DOUBLE: got %q", env["DOUBLE"])
	}
}

func TestParseFile_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 {
		t.Errorf("expected 1 key, got %d", len(env))
	}
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export MY_VAR=123\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["MY_VAR"] != "123" {
		t.Errorf("MY_VAR: got %q", env["MY_VAR"])
	}
}

func TestParseFile_MissingEquals(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
