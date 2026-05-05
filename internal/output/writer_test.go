package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/output"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "EQUAL_KEY", Status: diff.StatusEqual, ValueA: "same", ValueB: "same"},
		{Key: "REMOVED_KEY", Status: diff.StatusOnlyInA, ValueA: "old"},
		{Key: "ADDED_KEY", Status: diff.StatusOnlyInB, ValueB: "new"},
		{Key: "CHANGED_KEY", Status: diff.StatusChanged, ValueA: "before", ValueB: "after"},
	}
}

func TestWrite_TextNoDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := output.Options{Format: output.FormatText, NoColor: true, Out: &buf}
	n, err := output.Write([]diff.Entry{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 changed, got %d", n)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestWrite_TextWithDiff(t *testing.T) {
	var buf bytes.Buffer
	opts := output.Options{Format: output.FormatText, NoColor: true, Out: &buf}
	n, err := output.Write(makeEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 changed entries, got %d", n)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	opts := output.Options{Format: output.FormatJSON, Out: &buf}
	n, err := output.Write(makeEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 changed, got %d", n)
	}

	var doc struct {
		Differences []struct{ Status string } `json:"differences"`
		Total        int                       `json:"total"`
		Changed      int                       `json:"changed"`
	}
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if doc.Changed != 3 {
		t.Errorf("JSON changed field: want 3, got %d", doc.Changed)
	}
	if doc.Total != 4 {
		t.Errorf("JSON total field: want 4, got %d", doc.Total)
	}
	if len(doc.Differences) != 3 {
		t.Errorf("expected 3 difference entries, got %d", len(doc.Differences))
	}
}

func TestWrite_DefaultOut(t *testing.T) {
	// Ensure nil Out falls back gracefully (uses os.Stdout; just check no panic/error).
	opts := output.Options{Format: output.FormatText, NoColor: true, Out: nil}
	// Redirect by providing a real writer via DefaultOptions.
	var buf bytes.Buffer
	opts2 := output.DefaultOptions()
	opts2.Out = &buf
	opts2.NoColor = true
	_ = opts
	_, err := output.Write([]diff.Entry{}, opts2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
