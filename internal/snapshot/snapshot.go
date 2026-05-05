// Package snapshot provides functionality to capture and persist
// environment variable sets to disk for later comparison.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a saved environment variable set with metadata.
type Snapshot struct {
	Label     string            `json:"label"`
	Source    string            `json:"source"`
	CapturedAt time.Time        `json:"captured_at"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, label, source string, env map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("snapshot: create directory: %w", err)
	}

	snap := Snapshot{
		Label:      label,
		Source:     source,
		CapturedAt: time.Now().UTC(),
		Env:        env,
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}
