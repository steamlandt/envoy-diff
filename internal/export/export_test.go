package export_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/export"
)

func entries() []diff.Entry {
	return []diff.Entry{
		{Key: "APP_ENV", ValueA: "dev", ValueB: "prod", Status: diff.StatusChanged},
		{Key: "DEBUG", ValueA: "", ValueB: "true", Status: diff.StatusAdded},
		{Key: "OLD_KEY", ValueA: "x", ValueB: "", Status: diff.StatusRemoved},
		{Key: "PORT", ValueA: "8080", ValueB: "8080", Status: diff.StatusUnchanged},
	}
}

func TestWrite_ShellDefault(t *testing.T) {
	var buf strings.Builder
	err := export.Write(&buf, entries(), export.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=prod") {
		t.Errorf("expected export APP_ENV=prod, got:\n%s", out)
	}
	if !strings.Contains(out, "export PORT=8080") {
		t.Errorf("expected PORT in output, got:\n%s", out)
	}
	if strings.Contains(out, "OLD_KEY") {
		t.Errorf("removed key should be omitted, got:\n%s", out)
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	var buf strings.Builder
	err := export.Write(&buf, entries(), export.Options{Format: export.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "export ") {
		t.Errorf("dotenv format should not contain 'export', got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV=prod") {
		t.Errorf("expected APP_ENV=prod, got:\n%s", out)
	}
}

func TestWrite_OnlyChanged(t *testing.T) {
	var buf strings.Builder
	err := export.Write(&buf, entries(), export.Options{OnlyChanged: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "PORT") {
		t.Errorf("unchanged key PORT should be omitted, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("changed key APP_ENV should be present, got:\n%s", out)
	}
	if !strings.Contains(out, "DEBUG") {
		t.Errorf("added key DEBUG should be present, got:\n%s", out)
	}
}

func TestWrite_QuotesSpecialChars(t *testing.T) {
	special := []diff.Entry{
		{Key: "MSG", ValueA: "", ValueB: "hello world", Status: diff.StatusAdded},
	}
	var buf strings.Builder
	if err := export.Write(&buf, special, export.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "'hello world'") {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestWrite_EmptyValue(t *testing.T) {
	e := []diff.Entry{
		{Key: "EMPTY", ValueA: "old", ValueB: "", Status: diff.StatusChanged},
	}
	var buf strings.Builder
	if err := export.Write(&buf, e, export.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `EMPTY=""`) {
		t.Errorf("expected empty quoted value, got: %s", out)
	}
}
