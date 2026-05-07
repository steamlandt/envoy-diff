package normalize

import (
	"testing"
)

func TestApply_NoOp(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := Options{}
	r := Apply(env, opts)

	if len(r.Env) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.Env))
	}
	if r.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", r.Env["FOO"])
	}
	if len(r.Renamed) != 0 {
		t.Errorf("expected no renamed keys, got %v", r.Renamed)
	}
	if len(r.Dropped) != 0 {
		t.Errorf("expected no dropped keys, got %v", r.Dropped)
	}
}

func TestApply_TrimSpace(t *testing.T) {
	env := map[string]string{"  FOO  ": "  bar  "}
	opts := Options{TrimSpace: true}
	r := Apply(env, opts)

	if _, ok := r.Env["FOO"]; !ok {
		t.Fatalf("expected key FOO after trimming, got %v", r.Env)
	}
	if r.Env["FOO"] != "bar" {
		t.Errorf("expected value 'bar' after trimming, got %q", r.Env["FOO"])
	}
}

func TestApply_UpperKeys(t *testing.T) {
	env := map[string]string{"foo": "1", "Bar": "2", "BAZ": "3"}
	opts := Options{UpperKeys: true}
	r := Apply(env, opts)

	for _, k := range []string{"FOO", "BAR", "BAZ"} {
		if _, ok := r.Env[k]; !ok {
			t.Errorf("expected key %s in result", k)
		}
	}
	// foo and Bar were renamed
	if len(r.Renamed) != 2 {
		t.Errorf("expected 2 renamed keys, got %d: %v", len(r.Renamed), r.Renamed)
	}
}

func TestApply_DeduplicateKeys(t *testing.T) {
	// Simulate duplicate by using UpperKeys so two distinct originals collapse.
	env := map[string]string{"foo": "first", "FOO": "second"}
	opts := Options{UpperKeys: true, DeduplicateKeys: true}
	r := Apply(env, opts)

	if len(r.Env) != 1 {
		t.Fatalf("expected 1 entry after dedup, got %d", len(r.Env))
	}
	if len(r.Dropped) != 1 {
		t.Errorf("expected 1 dropped key, got %d: %v", len(r.Dropped), r.Dropped)
	}
}

func TestApply_DefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace=true in defaults")
	}
	if opts.UpperKeys {
		t.Error("expected UpperKeys=false in defaults")
	}
	if !opts.DeduplicateKeys {
		t.Error("expected DeduplicateKeys=true in defaults")
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	r := Apply(map[string]string{}, DefaultOptions())
	if len(r.Env) != 0 {
		t.Errorf("expected empty result, got %v", r.Env)
	}
}
