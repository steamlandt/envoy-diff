package diff

import (
	"testing"
)

func TestIntersect_NoCommonKeys(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	got := Intersect(a, b, DefaultIntersectOptions())
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestIntersect_AllSame(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "1", "BAR": "2"}
	got := Intersect(a, b, DefaultIntersectOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Changed {
			t.Errorf("key %q should not be marked changed", r.Key)
		}
	}
}

func TestIntersect_SomeChanged(t *testing.T) {
	a := map[string]string{"FOO": "old", "BAR": "same"}
	b := map[string]string{"FOO": "new", "BAR": "same"}
	got := Intersect(a, b, DefaultIntersectOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	// results must be sorted: BAR before FOO
	if got[0].Key != "BAR" || got[1].Key != "FOO" {
		t.Fatalf("unexpected order: %v", got)
	}
	if got[0].Changed {
		t.Errorf("BAR should not be changed")
	}
	if !got[1].Changed {
		t.Errorf("FOO should be changed")
	}
}

func TestIntersect_OnlyChanged(t *testing.T) {
	a := map[string]string{"FOO": "old", "BAR": "same", "BAZ": "x"}
	b := map[string]string{"FOO": "new", "BAR": "same", "BAZ": "x"}
	opts := IntersectOptions{OnlyChanged: true}
	got := Intersect(a, b, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d: %v", len(got), got)
	}
	if got[0].Key != "FOO" {
		t.Errorf("expected FOO, got %q", got[0].Key)
	}
	if got[0].ValueA != "old" || got[0].ValueB != "new" {
		t.Errorf("unexpected values: %+v", got[0])
	}
}

func TestIntersect_PartialOverlap(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"B": "2", "C": "99", "D": "4"}
	got := Intersect(a, b, DefaultIntersectOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 results (B,C), got %d: %v", len(got), got)
	}
	if got[0].Key != "B" || got[1].Key != "C" {
		t.Fatalf("unexpected keys: %v", got)
	}
	if got[1].Changed != true {
		t.Errorf("C should be marked changed")
	}
}
