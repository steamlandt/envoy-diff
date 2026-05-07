package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/truncate"
)

func TestValue_ShortValue_Unchanged(t *testing.T) {
	opts := truncate.DefaultOptions()
	got := truncate.Value("hello", opts)
	if got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestValue_ExactMaxLen_Unchanged(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 5
	opts.ShowLength = false
	s := "abcde"
	got := truncate.Value(s, opts)
	if got != s {
		t.Fatalf("expected %q, got %q", s, got)
	}
}

func TestValue_LongValue_Truncated(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 5
	opts.ShowLength = false
	s := "abcdefghij"
	got := truncate.Value(s, opts)
	if !strings.HasPrefix(got, "abcde") {
		t.Fatalf("expected prefix %q in %q", "abcde", got)
	}
	if !strings.Contains(got, opts.Suffix) {
		t.Fatalf("expected suffix %q in %q", opts.Suffix, got)
	}
}

func TestValue_ShowLength(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 3
	opts.ShowLength = true
	s := "abcdefgh" // 8 runes
	got := truncate.Value(s, opts)
	if !strings.Contains(got, "[8]") {
		t.Fatalf("expected length annotation [8] in %q", got)
	}
}

func TestValue_CustomSuffix(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 4
	opts.Suffix = "~~"
	opts.ShowLength = false
	s := "hello world"
	got := truncate.Value(s, opts)
	if !strings.HasSuffix(got, "~~") {
		t.Fatalf("expected suffix ~~ in %q", got)
	}
}

func TestValue_Unicode(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 3
	opts.ShowLength = false
	s := "日本語テスト" // 6 runes
	got := truncate.Value(s, opts)
	if !strings.HasPrefix(got, "日本語") {
		t.Fatalf("expected rune-based prefix in %q", got)
	}
}

func TestMap_TruncatesValues(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 5
	opts.ShowLength = false
	env := map[string]string{
		"SHORT": "hi",
		"LONG":  "this is a very long value",
	}
	out := truncate.Map(env, opts)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged, got %q", out["SHORT"])
	}
	if len([]rune(out["LONG"])) <= 5 && !strings.Contains(out["LONG"], opts.Suffix) {
		t.Errorf("LONG should be truncated, got %q", out["LONG"])
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	opts := truncate.DefaultOptions()
	opts.MaxLen = 3
	env := map[string]string{"KEY": "abcdefgh"}
	original := env["KEY"]
	truncate.Map(env, opts)
	if env["KEY"] != original {
		t.Errorf("input map was mutated")
	}
}
