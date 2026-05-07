package scope_test

import (
	"testing"

	"github.com/yourorg/envoy-diff/internal/scope"
)

func TestInject_Basic(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	opts := scope.DefaultOptions()
	got := scope.Inject("DB", env, opts)

	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestInject_LowercasePrefix_NormalisedToUpper(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	opts := scope.DefaultOptions()
	got := scope.Inject("app", env, opts)
	if _, ok := got["APP_KEY"]; !ok {
		t.Errorf("expected APP_KEY to exist, got %v", got)
	}
}

func TestInject_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"X": "1"}
	opts := scope.DefaultOptions()
	scope.Inject("NS", env, opts)
	if _, ok := env["NS_X"]; ok {
		t.Error("Inject must not mutate the input map")
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "myapp",
	}
	opts := scope.DefaultOptions()
	got := scope.Extract("DB", env, opts)

	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", got["HOST"])
	}
	if got["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", got["PORT"])
	}
	if _, ok := got["APP_NAME"]; ok {
		t.Error("APP_NAME should not appear in DB scope")
	}
}

func TestExtract_KeepPrefix(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := scope.DefaultOptions()
	opts.StripOnExtract = false
	got := scope.Extract("DB", env, opts)
	if _, ok := got["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST to be retained, got %v", got)
	}
}

func TestExtract_EmptyResult(t *testing.T) {
	env := map[string]string{"APP_NAME": "x"}
	opts := scope.DefaultOptions()
	got := scope.Extract("DB", env, opts)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestInjectExtract_RoundTrip(t *testing.T) {
	original := map[string]string{"USER": "alice", "PASS": "secret"}
	opts := scope.DefaultOptions()
	injected := scope.Inject("AUTH", original, opts)
	recovered := scope.Extract("AUTH", injected, opts)

	for k, v := range original {
		if recovered[k] != v {
			t.Errorf("round-trip mismatch for %s: want %q got %q", k, v, recovered[k])
		}
	}
}
