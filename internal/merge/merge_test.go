package merge

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Env["FOO"] != "1" || res.Env["BAR"] != "2" || res.Env["BAZ"] != "3" {
		t.Errorf("unexpected env map: %v", res.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := Merge(StrategyFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := Merge(StrategyLast, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "from-b" {
		t.Errorf("expected 'from-b', got %q", res.Env["KEY"])
	}
}

func TestMerge_StrategyError(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	_, err := Merge(StrategyError, a, b)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestMerge_StrategyError_NoConflict(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	res, err := Merge(StrategyError, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := Merge(StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}

func TestMerge_ConflictValues(t *testing.T) {
	a := map[string]string{"X": "alpha"}
	b := map[string]string{"X": "beta"}

	res, _ := Merge(StrategyFirst, a, b)
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict")
	}
	c := res.Conflicts[0]
	if c.Key != "X" {
		t.Errorf("expected conflict key 'X', got %q", c.Key)
	}
	if len(c.Values) != 2 {
		t.Errorf("expected 2 conflict values, got %d", len(c.Values))
	}
}
