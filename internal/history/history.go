// Package history manages a log of diff runs, storing results keyed by
// timestamp so users can review how their environment changed over time.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/example/envoy-diff/internal/diff"
)

// Entry is a single recorded diff run.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	SourceA   string           `json:"source_a"`
	SourceB   string           `json:"source_b"`
	Results   []diff.Entry     `json:"results"`
}

// Record appends a new entry to the history file stored under dir.
// The file is named history.json and entries are kept in chronological order.
func Record(dir string, e Entry) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("history: create dir: %w", err)
	}

	entries, _ := load(dir) // ignore error — start fresh if missing/corrupt
	entries = append(entries, e)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "history.json"), data, 0o644)
}

// List returns all recorded entries from dir, sorted oldest-first.
func List(dir string) ([]Entry, error) {
	entries, err := load(dir)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

// Prune removes entries older than maxAge from the history file.
func Prune(dir string, maxAge time.Duration) error {
	entries, err := load(dir)
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxAge)
	filtered := entries[:0]
	for _, e := range entries {
		if e.Timestamp.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "history.json"), data, 0o644)
}

func load(dir string) ([]Entry, error) {
	data, err := os.ReadFile(filepath.Join(dir, "history.json"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("history: unmarshal: %w", err)
	}
	return entries, nil
}
