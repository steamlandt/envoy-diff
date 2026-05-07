package diff

import "testing"

func TestApply_Unchanged(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	entries := []Entry{{Key: "FOO", ValueA: "bar", ValueB: "bar", Status: StatusUnchanged}}
	out, err := Apply(base, entries, DefaultApplyOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected bar, got %q", out["FOO"])
	}
}

func TestApply_Added(t *testing.T) {
	base := map[string]string{}
	entries := []Entry{{Key: "NEW", ValueB: "hello", Status: StatusAdded}}
	out, err := Apply(base, entries, DefaultApplyOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["NEW"] != "hello" {
		t.Errorf("expected hello, got %q", out["NEW"])
	}
}

func TestApply_Removed(t *testing.T) {
	base := map[string]string{"OLD": "gone"}
	entries := []Entry{{Key: "OLD", ValueA: "gone", Status: StatusRemoved}}
	out, err := Apply(base, entries, DefaultApplyOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
}

func TestApply_SkipRemoved(t *testing.T) {
	base := map[string]string{"OLD": "keep"}
	entries := []Entry{{Key: "OLD", ValueA: "keep", Status: StatusRemoved}}
	opts := ApplyOptions{SkipRemoved: true}
	out, err := Apply(base, entries, opts)
	if err != nil {
		t.Fatal(err)
	}
	if out["OLD"] != "keep" {
		t.Errorf("expected OLD to be kept, got %q", out["OLD"])
	}
}

func TestApply_Changed(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	entries := []Entry{{Key: "KEY", ValueA: "old", ValueB: "new", Status: StatusChanged}}
	out, err := Apply(base, entries, DefaultApplyOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "new" {
		t.Errorf("expected new, got %q", out["KEY"])
	}
}

func TestApply_FailOnConflict(t *testing.T) {
	base := map[string]string{"KEY": "diverged"}
	entries := []Entry{{Key: "KEY", ValueA: "old", ValueB: "new", Status: StatusChanged}}
	opts := ApplyOptions{FailOnConflict: true}
	_, err := Apply(base, entries, opts)
	if err == nil {
		t.Error("expected conflict error")
	}
}

func TestApply_NoConflict_WhenValuesMatch(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	entries := []Entry{{Key: "KEY", ValueA: "old", ValueB: "new", Status: StatusChanged}}
	opts := ApplyOptions{FailOnConflict: true}
	out, err := Apply(base, entries, opts)
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "new" {
		t.Errorf("expected new, got %q", out["KEY"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	entries := []Entry{{Key: "A", ValueA: "1", ValueB: "2", Status: StatusChanged}}
	_, err := Apply(base, entries, DefaultApplyOptions())
	if err != nil {
		t.Fatal(err)
	}
	if base["A"] != "1" {
		t.Error("Apply must not mutate the base map")
	}
}
