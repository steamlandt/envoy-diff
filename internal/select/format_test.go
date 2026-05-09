package sel_test

import (
	"bytes"
	"strings"
	"testing"

	sel "github.com/yourorg/envoy-diff/internal/select"
	"github.com/yourorg/envoy-diff/internal/diff"
)

func TestFormatText_NoSelection(t *testing.T) {
	var buf bytes.Buffer
	sel.FormatText(&buf, nil, 3, false)
	if !strings.Contains(buf.String(), "no entries selected") {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestFormatText_ShowsKeys(t *testing.T) {
	entries := []diff.Entry{
		{Key: "FOO", ValueA: "1", ValueB: "1"},
		{Key: "BAR", ValueA: "2", ValueB: "2"},
	}
	var buf bytes.Buffer
	sel.FormatText(&buf, entries, 5, false)
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %q", out)
	}
	if !strings.Contains(out, "BAR") {
		t.Errorf("expected BAR in output, got: %q", out)
	}
}

func TestFormatText_Summary(t *testing.T) {
	entries := []diff.Entry{
		{Key: "FOO", ValueA: "1", ValueB: "1"},
	}
	var buf bytes.Buffer
	sel.FormatText(&buf, entries, 4, false)
	out := buf.String()
	if !strings.Contains(out, "1 selected") {
		t.Errorf("expected '1 selected' in output, got: %q", out)
	}
	if !strings.Contains(out, "3 dropped") {
		t.Errorf("expected '3 dropped' in output, got: %q", out)
	}
}

func TestFormatText_ColorCodes(t *testing.T) {
	entries := []diff.Entry{{Key: "X", ValueA: "v", ValueB: "v"}}
	var buf bytes.Buffer
	sel.FormatText(&buf, entries, 1, true)
	if !strings.Contains(buf.String(), "\033[") {
		t.Error("expected ANSI codes in colour output")
	}
}

func TestParseKeys_Basic(t *testing.T) {
	got := sel.ParseKeys("FOO, BAR, BAZ")
	if len(got) != 3 || got[0] != "FOO" || got[1] != "BAR" || got[2] != "BAZ" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestParseKeys_Empty(t *testing.T) {
	if got := sel.ParseKeys(""); len(got) != 0 {
		t.Fatalf("expected nil/empty, got %v", got)
	}
}

func TestParseKeys_SkipsBlanks(t *testing.T) {
	got := sel.ParseKeys("FOO,,BAR")
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %v", got)
	}
}
