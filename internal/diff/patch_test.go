package diff

import (
	"strings"
	"testing"
)

func TestPatch_Unified_Empty(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", ValueA: "bar", ValueB: "bar", Status: StatusUnchanged},
	}
	out := Patch(entries, PatchFormatUnified)
	if out != "" {
		t.Errorf("expected empty patch, got %q", out)
	}
}

func TestPatch_Unified_Added(t *testing.T) {
	entries := []Entry{
		{Key: "NEW", ValueA: "", ValueB: "hello", Status: StatusAdded},
	}
	out := Patch(entries, PatchFormatUnified)
	if !strings.Contains(out, "+ NEW=hello") {
		t.Errorf("expected added line, got %q", out)
	}
}

func TestPatch_Unified_Removed(t *testing.T) {
	entries := []Entry{
		{Key: "OLD", ValueA: "gone", ValueB: "", Status: StatusRemoved},
	}
	out := Patch(entries, PatchFormatUnified)
	if !strings.Contains(out, "- OLD=gone") {
		t.Errorf("expected removed line, got %q", out)
	}
}

func TestPatch_Unified_Changed(t *testing.T) {
	entries := []Entry{
		{Key: "KEY", ValueA: "old", ValueB: "new", Status: StatusChanged},
	}
	out := Patch(entries, PatchFormatUnified)
	if !strings.Contains(out, "- KEY=old") || !strings.Contains(out, "+ KEY=new") {
		t.Errorf("expected both diff lines, got %q", out)
	}
}

func TestPatch_Shell_AddedAndRemoved(t *testing.T) {
	entries := []Entry{
		{Key: "A", ValueA: "", ValueB: "1", Status: StatusAdded},
		{Key: "B", ValueA: "2", ValueB: "", Status: StatusRemoved},
	}
	out := Patch(entries, PatchFormatShell)
	if !strings.Contains(out, "export A=") {
		t.Errorf("expected export line for A, got %q", out)
	}
	if !strings.Contains(out, "unset B") {
		t.Errorf("expected unset line for B, got %q", out)
	}
}

func TestPatch_Dotenv_RemovedIsComment(t *testing.T) {
	entries := []Entry{
		{Key: "GONE", ValueA: "x", ValueB: "", Status: StatusRemoved},
	}
	out := Patch(entries, PatchFormatDotenv)
	if !strings.Contains(out, "# removed: GONE") {
		t.Errorf("expected comment for removed key, got %q", out)
	}
}

func TestPatch_SortedOutput(t *testing.T) {
	entries := []Entry{
		{Key: "ZZZ", ValueA: "", ValueB: "z", Status: StatusAdded},
		{Key: "AAA", ValueA: "", ValueB: "a", Status: StatusAdded},
	}
	out := Patch(entries, PatchFormatUnified)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "+ AAA") {
		t.Errorf("expected AAA first, got %q", lines[0])
	}
}
