// Package source provides a unified interface for loading environment
// variables from different origins: .env files or running processes.
package source

import (
	"fmt"

	"github.com/yourorg/envoy-diff/internal/parser"
	"github.com/yourorg/envoy-diff/internal/process"
)

// Type represents the kind of environment source.
type Type string

const (
	// File indicates a .env file source.
	File Type = "file"
	// PID indicates a running process source.
	PID Type = "pid"
	// Self indicates the current process environment.
	Self Type = "self"
)

// Source describes an environment variable origin.
type Source struct {
	Kind Type
	// Path is set when Kind == File.
	Path string
	// PID is set when Kind == PID.
	PID int
}

// Load reads environment variables from the source and returns them as a
// map of key → value pairs.
func Load(s Source) (map[string]string, error) {
	switch s.Kind {
	case File:
		if s.Path == "" {
			return nil, fmt.Errorf("source: file path must not be empty")
		}
		return parser.ParseFile(s.Path)
	case PID:
		if s.PID <= 0 {
			return nil, fmt.Errorf("source: PID must be a positive integer, got %d", s.PID)
		}
		return process.ReadPID(s.PID)
	case Self:
		return process.ReadSelf(), nil
	default:
		return nil, fmt.Errorf("source: unknown source kind %q", s.Kind)
	}
}

// String returns a human-readable label for the source, suitable for diff
// output headers.
func (s Source) String() string {
	switch s.Kind {
	case File:
		return fmt.Sprintf("file:%s", s.Path)
	case PID:
		return fmt.Sprintf("pid:%d", s.PID)
	case Self:
		return "self"
	default:
		return fmt.Sprintf("unknown(%s)", s.Kind)
	}
}
