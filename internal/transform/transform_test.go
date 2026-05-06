package transform_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/transform"
)

func baseEnv() map[string]string {
	return map[string]string{
		"db_host": "localhost",
		"db_port": "5432",
		"API_KEY":  "secret",
	}
}

func TestApply_NoOp(t *testing.T) {
	out, results, err := transform.Apply(baseEnv(), transform.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
	for _, r := range results {
		if r.Skipped {
			t.Errorf("unexpected skip for key %s", r.OriginalKey)
		}
	}
}

func TestApply_UpperKeys(t *testing.T) {
	opts := transform.DefaultOptions()
	opts.UpperKeys = true
	out, _, err := transform.Apply(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k := range out {
		for _, c := range k {
			if c >= 'a' && c <= 'z' {
				t.Errorf("key %q still contains lowercase", k)
			}
		}
	}
}

func TestApply_LowerKeys(t *testing.T) {
	opts := transform.DefaultOptions()
	opts.LowerKeys = true
	out, _, err := transform.Apply(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["api_key"]; !ok {
		t.Error("expected api_key in output")
	}
}

func TestApply_MutuallyExclusiveCase(t *testing.T) {
	opts := transform.DefaultOptions()
	opts.UpperKeys = true
	opts.LowerKeys = true
	_, _, err := transform.Apply(baseEnv(), opts)
	if err == nil {
		t.Error("expected error for conflicting case options")
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	opts := transform.DefaultOptions()
	opts.KeyPrefix = "APP_"
	out, _, err := transform.Apply(map[string]string{"HOST": "localhost"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["APP_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %v", out)
	}
}

func TestApply_TrimValues(t *testing.T) {
	opts := transform.DefaultOptions()
	opts.TrimValues = true
	out, _, err := transform.Apply(map[string]string{"KEY": "  value  "}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApply_CollisionSkipped(t *testing.T) {
	// lower-casing both DB_HOST and db_host should cause a collision
	opts := transform.DefaultOptions()
	opts.LowerKeys = true
	env := map[string]string{"DB_HOST": "a", "db_host": "b"}
	_, results, err := transform.Apply(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	skipped := 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		}
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped result, got %d", skipped)
	}
}
