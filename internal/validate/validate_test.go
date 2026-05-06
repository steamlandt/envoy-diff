package validate_test

import (
	"testing"

	"github.com/your-org/envoy-diff/internal/validate"
)

func TestApply_NoIssues(t *testing.T) {
	env := map[string]string{
		"FOO":  "bar",
		"_BAZ": "qux",
		"MY_VAR_1": "hello",
	}
	issues := validate.Apply(env, validate.DefaultOptions())
	for _, i := range issues {
		if i.Severity == "error" {
			t.Errorf("unexpected error: %s", i)
		}
	}
}

func TestApply_EmptyKey(t *testing.T) {
	env := map[string]string{"": "value"}
	issues := validate.Apply(env, validate.DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected issue for empty key")
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestApply_InvalidChars(t *testing.T) {
	env := map[string]string{"MY-VAR": "val"}
	issues := validate.Apply(env, validate.DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected issue for invalid key character")
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error, got %s", issues[0].Severity)
	}
}

func TestApply_EmptyValue_Warning(t *testing.T) {
	env := map[string]string{"FOO": ""}
	opts := validate.DefaultOptions()
	opts.WarnEmptyValue = true
	issues := validate.Apply(env, opts)
	if len(issues) == 0 {
		t.Fatal("expected warning for empty value")
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestApply_EmptyValue_Suppressed(t *testing.T) {
	env := map[string]string{"FOO": ""}
	opts := validate.DefaultOptions()
	opts.WarnEmptyValue = false
	issues := validate.Apply(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestApply_ValueTooLong(t *testing.T) {
	env := map[string]string{"BIG": string(make([]byte, 100))}
	opts := validate.DefaultOptions()
	opts.MaxValueLen = 10
	issues := validate.Apply(env, opts)
	if len(issues) == 0 {
		t.Fatal("expected issue for oversized value")
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error, got %s", issues[0].Severity)
	}
}

func TestApply_DotInKey_Allowed(t *testing.T) {
	env := map[string]string{"my.var": "val"}
	opts := validate.DefaultOptions()
	opts.AllowDotInKey = true
	issues := validate.Apply(env, opts)
	for _, i := range issues {
		if i.Severity == "error" {
			t.Errorf("unexpected error: %s", i)
		}
	}
}

func TestHasErrors(t *testing.T) {
	issues := []validate.Issue{
		{Key: "X", Message: "bad", Severity: "warning"},
		{Key: "Y", Message: "worse", Severity: "error"},
	}
	if !validate.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}

	warningsOnly := []validate.Issue{
		{Key: "X", Message: "hmm", Severity: "warning"},
	}
	if validate.HasErrors(warningsOnly) {
		t.Error("expected HasErrors to return false for warnings only")
	}
}
