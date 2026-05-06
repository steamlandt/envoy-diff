package rename_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envoy-diff/internal/rename"
)

func TestApply_NoMapping(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := rename.Apply(env, nil, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestApply_SimpleRename(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	out, err := rename.Apply(env, map[string]string{"OLD_KEY": "NEW_KEY"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("old key should have been removed")
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
}

func TestApply_MissingSource_Ignored(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	_, err := rename.Apply(env, map[string]string{"MISSING": "TARGET"}, rename.DefaultOptions())
	if err != nil {
		t.Fatalf("expected no error for missing key with FailOnMissing=false, got %v", err)
	}
}

func TestApply_MissingSource_FailOnMissing(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := rename.DefaultOptions()
	opts.FailOnMissing = true
	_, err := rename.Apply(env, map[string]string{"MISSING": "TARGET"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key with FailOnMissing=true")
	}
}

func TestApply_TargetExists_NoOverwrite(t *testing.T) {
	env := map[string]string{"SRC": "v1", "DST": "v2"}
	_, err := rename.Apply(env, map[string]string{"SRC": "DST"}, rename.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when target exists and Overwrite=false")
	}
}

func TestApply_TargetExists_Overwrite(t *testing.T) {
	env := map[string]string{"SRC": "new", "DST": "old"}
	opts := rename.DefaultOptions()
	opts.Overwrite = true
	out, err := rename.Apply(env, map[string]string{"SRC": "DST"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DST"] != "new" {
		t.Errorf("expected DST=new, got %q", out["DST"])
	}
}

func TestApply_DuplicateTarget(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	_, err := rename.Apply(env, map[string]string{"A": "X", "B": "X"}, rename.DefaultOptions())
	if !errors.Is(err, rename.ErrDuplicateTarget) {
		t.Fatalf("expected ErrDuplicateTarget, got %v", err)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"OLD": "val"}
	_, _ = rename.Apply(env, map[string]string{"OLD": "NEW"}, rename.DefaultOptions())
	if _, ok := env["OLD"]; !ok {
		t.Error("Apply must not mutate the original map")
	}
}
