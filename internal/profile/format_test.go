package profile_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/profile"
)

func TestFormatList_Empty(t *testing.T) {
	var b strings.Builder
	profile.FormatList(&b, nil, false)
	if !strings.Contains(b.String(), "no profiles") {
		t.Errorf("expected 'no profiles' message, got %q", b.String())
	}
}

func TestFormatList_Names(t *testing.T) {
	var b strings.Builder
	profile.FormatList(&b, []string{"dev", "prod"}, false)
	out := b.String()
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Errorf("expected both names in output, got %q", out)
	}
}

func TestFormatList_ColorCodes(t *testing.T) {
	var b strings.Builder
	profile.FormatList(&b, []string{"staging"}, true)
	if !strings.Contains(b.String(), "\033[") {
		t.Error("expected ANSI color codes in color mode")
	}
}

func TestFormatProfile_ContainsKeys(t *testing.T) {
	var b strings.Builder
	p := profile.Profile{
		Name:    "dev",
		Env:     map[string]string{"FOO": "bar", "BAZ": "hello world"},
		SavedAt: time.Now(),
	}
	profile.FormatProfile(&b, p, false)
	out := b.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got %q", out)
	}
	if !strings.Contains(out, `BAZ="hello world"`) {
		t.Errorf("expected quoted value for BAZ, got %q", out)
	}
	if !strings.Contains(out, "# profile: dev") {
		t.Errorf("expected header comment, got %q", out)
	}
}

func TestFormatProfile_NoColor(t *testing.T) {
	var b strings.Builder
	p := profile.Profile{Name: "x", Env: map[string]string{"K": "v"}, SavedAt: time.Now()}
	profile.FormatProfile(&b, p, false)
	if strings.Contains(b.String(), "\033[") {
		t.Error("unexpected ANSI codes in no-color mode")
	}
}
