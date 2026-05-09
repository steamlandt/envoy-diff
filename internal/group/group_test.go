package group_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/group"
)

func TestApply_NoEntries(t *testing.T) {
	result := group.Apply(map[string]string{}, group.DefaultOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(result))
	}
}

func TestApply_SinglePrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	}
	result := group.Apply(env, group.DefaultOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}
	if result[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %q", result[0].Prefix)
	}
	if len(result[0].Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result[0].Entries))
	}
}

func TestApply_MultiplePrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"REDIS_HOST": "127.0.0.1",
		"REDIS_PORT": "6379",
	}
	result := group.Apply(env, group.DefaultOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}
	if result[0].Prefix != "DB" {
		t.Errorf("expected first group DB, got %q", result[0].Prefix)
	}
	if result[1].Prefix != "REDIS" {
		t.Errorf("expected second group REDIS, got %q", result[1].Prefix)
	}
}

func TestApply_NoSeparator_EmptyPrefix(t *testing.T) {
	env := map[string]string{
		"HOSTNAME": "myhost",
		"PATH":     "/usr/bin",
	}
	result := group.Apply(env, group.DefaultOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 group (empty prefix), got %d", len(result))
	}
	if result[0].Prefix != "" {
		t.Errorf("expected empty prefix, got %q", result[0].Prefix)
	}
}

func TestApply_MinGroupSize(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}
	opts := group.DefaultOptions()
	opts.MinGroupSize = 2
	result := group.Apply(env, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 group (DB only), got %d", len(result))
	}
	if result[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %q", result[0].Prefix)
	}
}

func TestApply_EntriesSortedByKey(t *testing.T) {
	env := map[string]string{
		"DB_PORT": "5432",
		"DB_HOST": "localhost",
		"DB_NAME": "mydb",
	}
	result := group.Apply(env, group.DefaultOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 group, got %d", len(result))
	}
	keys := make([]string, len(result[0].Entries))
	for i, e := range result[0].Entries {
		keys[i] = e.Key
	}
	expected := []string{"DB_HOST", "DB_NAME", "DB_PORT"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("entry[%d]: expected %q, got %q", i, k, keys[i])
		}
	}
}
