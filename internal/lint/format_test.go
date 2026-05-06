package lint

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	FormatText(&buf, nil, false)
	if !strings.Contains(buf.String(), "No lint issues") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestFormatText_ShowsError(t *testing.T) {
	var buf bytes.Buffer
	findings := []Finding{
		{Key: "FOO", Message: "duplicate", Severity: SeverityError},
	}
	FormatText(&buf, findings, false)
	out := buf.String()
	if !strings.Contains(out, "[error]") {
		t.Errorf("expected [error] label, got: %q", out)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key FOO in output, got: %q", out)
	}
	if !strings.Contains(out, "1 error") {
		t.Errorf("expected summary with 1 error, got: %q", out)
	}
}

func TestFormatText_ShowsWarning(t *testing.T) {
	var buf bytes.Buffer
	findings := []Finding{
		{Key: "bar", Message: "not uppercase", Severity: SeverityWarn},
	}
	FormatText(&buf, findings, false)
	out := buf.String()
	if !strings.Contains(out, "[warn]") {
		t.Errorf("expected [warn] label, got: %q", out)
	}
	if !strings.Contains(out, "1 warning") {
		t.Errorf("expected summary with 1 warning, got: %q", out)
	}
}

func TestFormatText_ErrorsBeforeWarnings(t *testing.T) {
	var buf bytes.Buffer
	findings := []Finding{
		{Key: "Z_WARN", Message: "warn msg", Severity: SeverityWarn},
		{Key: "A_ERR", Message: "err msg", Severity: SeverityError},
	}
	FormatText(&buf, findings, false)
	out := buf.String()
	errIdx := strings.Index(out, "[error]")
	warnIdx := strings.Index(out, "[warn]")
	if errIdx > warnIdx {
		t.Errorf("expected errors before warnings in output")
	}
}

func TestFormatText_ColorCodes(t *testing.T) {
	var buf bytes.Buffer
	findings := []Finding{
		{Key: "X", Message: "oops", Severity: SeverityError},
	}
	FormatText(&buf, findings, true)
	if !strings.Contains(buf.String(), ansiRed) {
		t.Error("expected ANSI red code in color output")
	}
}
