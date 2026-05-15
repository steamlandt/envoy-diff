// Package profile manages named environment profiles — saved collections
// of key/value pairs that can be loaded, listed, and deleted by name.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ErrNotFound is returned when a named profile does not exist.
var ErrNotFound = errors.New("profile not found")

// Profile holds a named snapshot of environment variables.
type Profile struct {
	Name    string            `json:"name"`
	Env     map[string]string `json:"env"`
	SavedAt time.Time         `json:"saved_at"`
}

// Save writes env to a named profile under dir.
func Save(dir, name string, env map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("profile: create dir: %w", err)
	}
	p := Profile{Name: name, Env: env, SavedAt: time.Now().UTC()}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("profile: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, name+".json"), data, 0o644)
}

// Load reads the named profile from dir.
func Load(dir, name string) (Profile, error) {
	data, err := os.ReadFile(filepath.Join(dir, name+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return Profile{}, ErrNotFound
	}
	if err != nil {
		return Profile{}, fmt.Errorf("profile: read: %w", err)
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("profile: unmarshal: %w", err)
	}
	return p, nil
}

// List returns the names of all saved profiles in dir, sorted.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("profile: list: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	sort.Strings(names)
	return names, nil
}

// Delete removes the named profile from dir.
func Delete(dir, name string) error {
	err := os.Remove(filepath.Join(dir, name+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}
	return err
}
