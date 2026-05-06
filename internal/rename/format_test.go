package rename_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/rename"
)

func TestFormatText_NoMapping(t *testing.T) {
	var sb strings.Builder
	rename.FormatText(&sb, nil, map[string]string{}, false)
	if !strings.Contains(sb.String(), "no renames") {
		t.Errorf("expected 'no renames' message, got: %q", sb.String())
	}
}

func TestFormatText_ShowsRename(t *testing.T) {
	var sb strings.Builder
	env := map[string]string{"OLD": "val"}
	mapping := map[string]string{"OLD": "NEW"}
	rename.FormatText(&sb, mapping, env, false)
	out := sb.String()
	if !strings.Contains(out, "OLD -> NEW") {
		t.Errorf("expected 'OLD -> NEW' in output, got: %q", out)
	}
}

func TestFormatText_ShowsSkipped(t *testing.T) {
	var sb strings.Builder
	env := map[string]string{}
	mapping := map[string]string{"MISSING": "TARGET"}
	rename.FormatText(&sb, mapping, env, false)
	out := sb.String()
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output, got: %q", out)
	}
}

func TestFormatText_ColorCodes(t *testing.T) {
	var sb strings.Builder
	env := map[string]string{"K": "v"}
	mapping := map[string]string{"K": "K2"}
	rename.FormatText(&sb, mapping, env, true)
	if !strings.Contains(sb.String(), "\033[") {
		t.Error("expected ANSI escape codes when color=true")
	}
}

func TestParseMapping_Basic(t *testing.T) {
	pairs := []string{"OLD=NEW", "FOO=BAR"}
	m := rename.ParseMapping(pairs)
	if m["OLD"] != "NEW" || m["FOO"] != "BAR" {
		t.Errorf("unexpected mapping: %v", m)
	}
}

func TestParseMapping_IgnoresInvalid(t *testing.T) {
	pairs := []string{"NOEQUALSSIGN", "=NODST", "NOSRC="}
	m := rename.ParseMapping(pairs)
	if len(m) != 0 {
		t.Errorf("expected empty mapping, got: %v", m)
	}
}

func TestParseMapping_Whitespace(t *testing.T) {
	pairs := []string{" OLD = NEW "}
	m := rename.ParseMapping(pairs)
	if m["OLD"] != "NEW" {
		t.Errorf("expected OLD->NEW after trimming, got: %v", m)
	}
}
