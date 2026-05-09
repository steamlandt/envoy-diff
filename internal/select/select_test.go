package sel_test

import (
	"testing"

	sel "github.com/yourorg/envoy-diff/internal/select"
	"github.com/yourorg/envoy-diff/internal/diff"
)

func makeEntries(pairs ...string) []diff.Entry {
	var out []diff.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, diff.Entry{Key: pairs[i], ValueA: pairs[i+1], ValueB: pairs[i+1]})
	}
	return out
}

func keys(entries []diff.Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NoFilter(t *testing.T) {
	in := makeEntries("FOO", "1", "BAR", "2", "BAZ", "3")
	out, err := sel.Apply(in, sel.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	in := makeEntries("FOO", "1", "BAR", "2", "BAZ", "3")
	out, err := sel.Apply(in, sel.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatal(err)
	}
	got := keys(out)
	if len(got) != 2 || got[0] != "FOO" || got[1] != "BAZ" {
		t.Fatalf("unexpected keys: %v", got)
	}
}

func TestApply_Exclude(t *testing.T) {
	in := makeEntries("FOO", "1", "BAR", "2", "BAZ", "3")
	out, err := sel.Apply(in, sel.Options{Exclude: []string{"BAR"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range out {
		if e.Key == "BAR" {
			t.Fatal("BAR should have been excluded")
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestApply_Pattern(t *testing.T) {
	in := makeEntries("DB_HOST", "h", "DB_PORT", "5432", "APP_NAME", "x")
	out, err := sel.Apply(in, sel.Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestApply_ExcludePattern(t *testing.T) {
	in := makeEntries("SECRET_KEY", "s", "DB_HOST", "h", "DB_PORT", "p")
	out, err := sel.Apply(in, sel.Options{ExcludePattern: "^DB_"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Key != "SECRET_KEY" {
		t.Fatalf("unexpected result: %v", keys(out))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	in := makeEntries("FOO", "1")
	_, err := sel.Apply(in, sel.Options{Pattern: "["})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_CaseInsensitiveKeys(t *testing.T) {
	in := makeEntries("foo", "1", "BAR", "2")
	out, err := sel.Apply(in, sel.Options{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Key != "foo" {
		t.Fatalf("expected foo entry, got %v", keys(out))
	}
}
