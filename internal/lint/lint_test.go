package lint

import (
	"testing"
)

func TestApply_NoIssues(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	findings := Apply(env, DefaultOptions())
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestApply_LowercaseKey(t *testing.T) {
	env := map[string]string{"host": "localhost"}
	findings := Apply(env, DefaultOptions())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityWarn {
		t.Errorf("expected warn, got %s", findings[0].Severity)
	}
}

func TestApply_DuplicateCaseInsensitive(t *testing.T) {
	env := map[string]string{
		"HOST": "a",
		"host": "b",
	}
	findings := Apply(env, DefaultOptions())
	// Should contain at least one duplicate error.
	var found bool
	for _, f := range findings {
		if f.Severity == SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one error finding for duplicate keys")
	}
}

func TestApply_ValueTooLong(t *testing.T) {
	long := make([]byte, 100)
	for i := range long {
		long[i] = 'x'
	}
	env := map[string]string{"BIG": string(long)}
	opts := DefaultOptions()
	opts.MaxValueLen = 10
	findings := Apply(env, opts)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityWarn {
		t.Errorf("expected warn, got %s", findings[0].Severity)
	}
}

func TestApply_DisabledChecks(t *testing.T) {
	env := map[string]string{"lower": "val"}
	opts := Options{CheckKeyUppercase: false}
	findings := Apply(env, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings with check disabled, got %d", len(findings))
	}
}

func TestHasErrors_True(t *testing.T) {
	findings := []Finding{{Key: "X", Severity: SeverityError}}
	if !HasErrors(findings) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	findings := []Finding{{Key: "X", Severity: SeverityWarn}}
	if HasErrors(findings) {
		t.Error("expected HasErrors to return false")
	}
}
