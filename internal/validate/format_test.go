package validate_test

import (
	"strings"
	"testing"

	"github.com/your-org/envoy-diff/internal/validate"
)

func TestFormatText_NoIssues(t *testing.T) {
	var sb strings.Builder
	validate.FormatText(&sb, nil, false)
	if !strings.Contains(sb.String(), "No validation issues") {
		t.Errorf("expected no-issues message, got: %s", sb.String())
	}
}

func TestFormatText_ShowsErrors(t *testing.T) {
	issues := []validate.Issue{
		{Key: "BAD KEY", Message: "contains whitespace", Severity: "error"},
	}
	var sb strings.Builder
	validate.FormatText(&sb, issues, false)
	out := sb.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR label in output, got: %s", out)
	}
	if !strings.Contains(out, "BAD KEY") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestFormatText_ShowsWarnings(t *testing.T) {
	issues := []validate.Issue{
		{Key: "EMPTY_VAL", Message: "value is empty", Severity: "warning"},
	}
	var sb strings.Builder
	validate.FormatText(&sb, issues, false)
	out := sb.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN label, got: %s", out)
	}
}

func TestFormatText_Summary(t *testing.T) {
	issues := []validate.Issue{
		{Key: "A", Message: "bad", Severity: "error"},
		{Key: "B", Message: "meh", Severity: "warning"},
	}
	var sb strings.Builder
	validate.FormatText(&sb, issues, false)
	out := sb.String()
	if !strings.Contains(out, "1 error(s), 1 warning(s)") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestFormatText_ErrorsBeforeWarnings(t *testing.T) {
	issues := []validate.Issue{
		{Key: "Z_WARN", Message: "empty", Severity: "warning"},
		{Key: "A_ERR", Message: "invalid", Severity: "error"},
	}
	var sb strings.Builder
	validate.FormatText(&sb, issues, false)
	out := sb.String()
	errIdx := strings.Index(out, "A_ERR")
	warnIdx := strings.Index(out, "Z_WARN")
	if errIdx > warnIdx {
		t.Error("expected errors to appear before warnings in output")
	}
}
