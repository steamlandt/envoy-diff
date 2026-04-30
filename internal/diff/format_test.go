package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormat_NoDiff(t *testing.T) {
	r := Result{
		OnlyInA:   map[string]string{},
		OnlyInB:   map[string]string{},
		Changed:   map[string][2]string{},
		Unchanged: map[string]string{"FOO": "bar"},
	}
	var buf bytes.Buffer
	Format(&buf, r, "A", "B")
	if !strings.Contains(buf.String(), "No differences") {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestFormat_OnlyInA(t *testing.T) {
	r := Result{
		OnlyInA:   map[string]string{"GONE": "val"},
		OnlyInB:   map[string]string{},
		Changed:   map[string][2]string{},
		Unchanged: map[string]string{},
	}
	var buf bytes.Buffer
	Format(&buf, r, "file1", "file2")
	out := buf.String()
	if !strings.Contains(out, "only in file1") {
		t.Errorf("expected label file1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "GONE=val") {
		t.Errorf("expected GONE=val in output, got:\n%s", out)
	}
}

func TestFormat_Changed(t *testing.T) {
	r := Result{
		OnlyInA:   map[string]string{},
		OnlyInB:   map[string]string{},
		Changed:   map[string][2]string{"PORT": {"8080", "9090"}},
		Unchanged: map[string]string{},
	}
	var buf bytes.Buffer
	Format(&buf, r, "dev", "prod")
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output")
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Errorf("expected both values in output, got:\n%s", out)
	}
}
