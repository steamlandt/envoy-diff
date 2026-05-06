package convert_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/convert"
)

func TestToSlice_Sorted(t *testing.T) {
	env := map[string]string{"Z": "26", "A": "1", "M": "13"}
	got := convert.ToSlice(env)
	want := []string{"A=1", "M=13", "Z=26"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestToSlice_Empty(t *testing.T) {
	got := convert.ToSlice(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestFromSlice_Basic(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux"}
	got := convert.FromSlice(pairs)
	if got["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", got["FOO"], "bar")
	}
	if got["BAZ"] != "qux" {
		t.Errorf("BAZ: got %q, want %q", got["BAZ"], "qux")
	}
}

func TestFromSlice_ValueWithEquals(t *testing.T) {
	got := convert.FromSlice([]string{"URL=http://x.com?a=1&b=2"})
	if got["URL"] != "http://x.com?a=1&b=2" {
		t.Errorf("unexpected value: %q", got["URL"])
	}
}

func TestFromSlice_NoEquals(t *testing.T) {
	got := convert.FromSlice([]string{"NOEQUALS"})
	v, ok := got["NOEQUALS"]
	if !ok {
		t.Fatal("expected key NOEQUALS to be present")
	}
	if v != "" {
		t.Errorf("expected empty value, got %q", v)
	}
}

func TestFromSlice_SkipsBlanks(t *testing.T) {
	got := convert.FromSlice([]string{"", "  ", "K=V"})
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
}

func TestMergeOverride(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	over := map[string]string{"B": "99", "C": "3"}
	got := convert.MergeOverride(base, over)
	if got["A"] != "1" {
		t.Errorf("A: got %q, want %q", got["A"], "1")
	}
	if got["B"] != "99" {
		t.Errorf("B: got %q, want %q", got["B"], "99")
	}
	if got["C"] != "3" {
		t.Errorf("C: got %q, want %q", got["C"], "3")
	}
	// ensure base is not mutated
	if base["B"] != "2" {
		t.Error("base map was mutated")
	}
}

func TestKeys_Sorted(t *testing.T) {
	env := map[string]string{"Z": "", "A": "", "M": ""}
	got := convert.Keys(env)
	want := []string{"A", "M", "Z"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q, want %q", i, got[i], want[i])
		}
	}
}
