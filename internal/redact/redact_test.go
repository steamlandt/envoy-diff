package redact_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/redact"
)

func TestApply_NoSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"HOME":  "/home/user",
		"SHELL": "/bin/bash",
		"PATH":  "/usr/bin:/bin",
	}
	out := redact.Apply(env, nil)
	for k, v := range env {
		if out[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestApply_DefaultPatterns(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD":  "s3cr3t",
		"API_KEY":      "abc123",
		"GITHUB_TOKEN": "ghp_xxx",
		"HOME":         "/home/user",
	}
	out := redact.Apply(env, nil)

	for _, sensitive := range []string{"DB_PASSWORD", "API_KEY", "GITHUB_TOKEN"} {
		if out[sensitive] != redact.Mask() {
			t.Errorf("key %s: expected redaction, got %q", sensitive, out[sensitive])
		}
	}
	if out["HOME"] != "/home/user" {
		t.Errorf("HOME should not be redacted, got %q", out["HOME"])
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_xxx",
		"APP_NAME":   "envoy-diff",
	}
	out := redact.Apply(env, []string{"STRIPE"})

	if out["STRIPE_KEY"] != redact.Mask() {
		t.Errorf("STRIPE_KEY should be redacted")
	}
	if out["APP_NAME"] != "envoy-diff" {
		t.Errorf("APP_NAME should not be redacted")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	env := map[string]string{
		"db_password": "lowercase",
		"Db_Password": "mixedcase",
	}
	out := redact.Apply(env, nil)
	for k := range env {
		if out[k] != redact.Mask() {
			t.Errorf("key %s should be redacted regardless of case", k)
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "original",
	}
	redact.Apply(env, nil)
	if env["SECRET_KEY"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	out := redact.Apply(map[string]string{}, nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
