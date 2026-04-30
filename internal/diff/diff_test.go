package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	if r.HasDifferences() {
		t.Fatal("expected no differences")
	}
	if len(r.Unchanged) != 2 {
		t.Fatalf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "EXTRA": "only"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	if _, ok := r.OnlyInA["EXTRA"]; !ok {
		t.Fatal("expected EXTRA in OnlyInA")
	}
	if len(r.OnlyInB) != 0 {
		t.Fatal("expected OnlyInB to be empty")
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "NEW": "val"}
	r := Compare(a, b)
	if _, ok := r.OnlyInB["NEW"]; !ok {
		t.Fatal("expected NEW in OnlyInB")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Compare(a, b)
	pair, ok := r.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Fatalf("unexpected changed values: %v", pair)
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "same"}
	b := map[string]string{"A": "changed", "D": "new", "C": "same"}
	r := Compare(a, b)
	if !r.HasDifferences() {
		t.Fatal("expected differences")
	}
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	if len(r.OnlyInA) != 1 {
		t.Fatalf("expected 1 only-in-A, got %d", len(r.OnlyInA))
	}
	if len(r.OnlyInB) != 1 {
		t.Fatalf("expected 1 only-in-B, got %d", len(r.OnlyInB))
	}
	if len(r.Unchanged) != 1 {
		t.Fatalf("expected 1 unchanged, got %d", len(r.Unchanged))
	}
}
