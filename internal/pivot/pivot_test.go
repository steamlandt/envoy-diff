package pivot_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envoy-diff/internal/pivot"
)

func TestBuild_KeyUnion(t *testing.T) {
	srcs := []string{"a", "b"}
	envs := map[string]map[string]string{
		"a": {"FOO": "1", "BAR": "x"},
		"b": {"FOO": "1", "BAZ": "z"},
	}
	table := pivot.Build(srcs, envs)
	if len(table.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(table.Rows))
	}
	if table.Rows[0].Key != "BAR" {
		t.Errorf("expected first key BAR, got %s", table.Rows[0].Key)
	}
}

func TestBuild_UniformDetection(t *testing.T) {
	srcs := []string{"a", "b"}
	envs := map[string]map[string]string{
		"a": {"FOO": "same", "BAR": "diff1"},
		"b": {"FOO": "same", "BAR": "diff2"},
	}
	table := pivot.Build(srcs, envs)
	for _, row := range table.Rows {
		switch row.Key {
		case "FOO":
			if !row.Uniform {
				t.Error("FOO should be uniform")
			}
		case "BAR":
			if row.Uniform {
				t.Error("BAR should not be uniform")
			}
		}
	}
}

func TestBuild_AbsentKey(t *testing.T) {
	srcs := []string{"a", "b"}
	envs := map[string]map[string]string{
		"a": {"ONLY_A": "val"},
		"b": {},
	}
	table := pivot.Build(srcs, envs)
	if len(table.Rows) != 1 {
		t.Fatalf("expected 1 row")
	}
	row := table.Rows[0]
	if row.Sources[0].Present != true {
		t.Error("source a should have key present")
	}
	if row.Sources[1].Present != false {
		t.Error("source b should not have key present")
	}
	if row.Uniform {
		t.Error("single-source key should not be uniform")
	}
}

func TestBuild_EmptySources(t *testing.T) {
	table := pivot.Build([]string{"a"}, map[string]map[string]string{"a": {}})
	if len(table.Rows) != 0 {
		t.Errorf("expected 0 rows for empty env")
	}
}

func TestFormatText_ShowsHeader(t *testing.T) {
	srcs := []string{"staging", "prod"}
	envs := map[string]map[string]string{
		"staging": {"PORT": "8080"},
		"prod":    {"PORT": "443"},
	}
	table := pivot.Build(srcs, envs)
	var buf bytes.Buffer
	pivot.FormatText(&buf, table, pivot.FormatOptions{Color: false})
	out := buf.String()
	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY header")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected staging header")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT row")
	}
}

func TestFormatText_OnlyDiffing(t *testing.T) {
	srcs := []string{"a", "b"}
	envs := map[string]map[string]string{
		"a": {"SAME": "v", "DIFF": "1"},
		"b": {"SAME": "v", "DIFF": "2"},
	}
	table := pivot.Build(srcs, envs)
	var buf bytes.Buffer
	pivot.FormatText(&buf, table, pivot.FormatOptions{Color: false, OnlyDiffing: true})
	out := buf.String()
	if strings.Contains(out, "SAME") {
		t.Error("SAME should be filtered out with OnlyDiffing")
	}
	if !strings.Contains(out, "DIFF") {
		t.Error("DIFF should appear with OnlyDiffing")
	}
}
