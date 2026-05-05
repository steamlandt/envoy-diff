package filter_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.local",
	"DB_PASSWORD": "secret",
	"LOG_LEVEL":   "info",
	"DEBUG":       "false",
}

func TestApply_NoFilter(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(sampleEnv) {
		t.Errorf("expected %d keys, got %d", len(sampleEnv), len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApply_Suffix(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", 2, out)
	}
}

func TestApply_Pattern(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := filter.Apply(sampleEnv, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid pattern, got nil")
	}
}

func TestApply_PrefixAndSuffix(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_", Suffix: "_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestApply_NoMatch(t *testing.T) {
	out, err := filter.Apply(sampleEnv, filter.Options{Prefix: "NOPE_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected 0 keys, got %d", len(out))
	}
}
