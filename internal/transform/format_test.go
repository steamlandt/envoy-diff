package transform_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/transform"
)

func TestFormatText_NoResults(t *testing.T) {
	var buf bytes.Buffer
	transform.FormatText(&buf, nil, false)
	if !strings.Contains(buf.String(), "No transformations") {
		t.Errorf("expected no-transformations message, got %q", buf.String())
	}
}

func TestFormatText_ShowsRename(t *testing.T) {
	results := []transform.Result{
		{OriginalKey: "db_host", NewKey: "DB_HOST", Value: "localhost"},
	}
	var buf bytes.Buffer
	transform.FormatText(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "db_host") || !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected rename line in output, got %q", out)
	}
	if !strings.Contains(out, "1 renamed") {
		t.Errorf("expected summary line, got %q", out)
	}
}

func TestFormatText_ShowsSkipped(t *testing.T) {
	results := []transform.Result{
		{OriginalKey: "DB_HOST", NewKey: "db_host", Skipped: true, Reason: "key collision after transformation"},
	}
	var buf bytes.Buffer
	transform.FormatText(&buf, results, false)
	out := buf.String()
	if !strings.Contains(out, "SKIP") {
		t.Errorf("expected SKIP label, got %q", out)
	}
	if !strings.Contains(out, "1 skipped") {
		t.Errorf("expected skipped summary, got %q", out)
	}
}

func TestFormatText_ColorCodes(t *testing.T) {
	results := []transform.Result{
		{OriginalKey: "foo", NewKey: "FOO", Value: "bar"},
	}
	var buf bytes.Buffer
	transform.FormatText(&buf, results, true)
	if !strings.Contains(buf.String(), "\033[") {
		t.Error("expected ANSI codes in color output")
	}
}

func TestFormatText_NoKeyChanges(t *testing.T) {
	// key unchanged, only value trimmed
	results := []transform.Result{
		{OriginalKey: "HOST", NewKey: "HOST", Value: "localhost"},
	}
	var buf bytes.Buffer
	transform.FormatText(&buf, results, false)
	if !strings.Contains(buf.String(), "No key changes") {
		t.Errorf("expected no-key-changes message, got %q", buf.String())
	}
}
