package resolve_test

import (
	"os"
	"testing"

	"github.com/yourorg/envoy-diff/internal/resolve"
)

func TestApply_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := resolve.Apply(env, resolve.DefaultOptions())
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApply_SimpleReference(t *testing.T) {
	env := map[string]string{
		"BASE": "/usr/local",
		"BIN":  "${BASE}/bin",
	}
	out := resolve.Apply(env, resolve.DefaultOptions())
	if got, want := out["BIN"], "/usr/local/bin"; got != want {
		t.Errorf("BIN = %q, want %q", got, want)
	}
}

func TestApply_ChainedReference(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out := resolve.Apply(env, resolve.DefaultOptions())
	if got, want := out["C"], "hello_world!"; got != want {
		t.Errorf("C = %q, want %q", got, want)
	}
}

func TestApply_FallbackToOS(t *testing.T) {
	os.Setenv("_RESOLVE_TEST_KEY", "from_os")
	t.Cleanup(func() { os.Unsetenv("_RESOLVE_TEST_KEY") })

	env := map[string]string{"VAL": "${_RESOLVE_TEST_KEY}"}
	opts := resolve.DefaultOptions()
	opts.FallbackToOS = true
	out := resolve.Apply(env, opts)
	if got, want := out["VAL"], "from_os"; got != want {
		t.Errorf("VAL = %q, want %q", got, want)
	}
}

func TestApply_NoFallbackToOS(t *testing.T) {
	os.Setenv("_RESOLVE_TEST_KEY2", "should_not_appear")
	t.Cleanup(func() { os.Unsetenv("_RESOLVE_TEST_KEY2") })

	env := map[string]string{"VAL": "${_RESOLVE_TEST_KEY2}"}
	opts := resolve.DefaultOptions()
	opts.FallbackToOS = false
	out := resolve.Apply(env, opts)
	if got := out["VAL"]; got != "" {
		t.Errorf("VAL = %q, want empty string", got)
	}
}

func TestApply_MaxDepthPreventsInfiniteLoop(t *testing.T) {
	// A -> ${A} would recurse forever without depth limit.
	env := map[string]string{"A": "${A}"}
	opts := resolve.Options{FallbackToOS: false, MaxDepth: 3}
	// Should not panic or hang.
	out := resolve.Apply(env, opts)
	_ = out["A"] // value is undefined but must not deadlock
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"X": "${Y}", "Y": "resolved"}
	orig := map[string]string{"X": "${Y}", "Y": "resolved"}
	resolve.Apply(env, resolve.DefaultOptions())
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("input mutated: env[%q] = %q, want %q", k, env[k], v)
		}
	}
}

func TestHasReferences(t *testing.T) {
	if resolve.HasReferences("no refs here") {
		t.Error("expected false for plain string")
	}
	if !resolve.HasReferences("${FOO}") {
		t.Error("expected true for ${FOO}")
	}
	if !resolve.HasReferences("$BAR") {
		t.Error("expected true for $BAR")
	}
}
