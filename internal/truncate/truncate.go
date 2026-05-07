// Package truncate provides utilities for truncating long environment variable
// values for display purposes, preserving a configurable prefix and suffix.
package truncate

import "fmt"

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		MaxLen:     80,
		Suffix:     "...",
		ShowLength: true,
	}
}

// Options controls how values are truncated.
type Options struct {
	// MaxLen is the maximum number of runes to show before truncating.
	MaxLen int
	// Suffix is appended when a value is truncated.
	Suffix string
	// ShowLength appends the total original length in brackets when truncated.
	ShowLength bool
}

// Value truncates s according to opts. If s is within MaxLen runes it is
// returned unchanged. Otherwise the first MaxLen runes are returned followed
// by Suffix and, optionally, the original length.
func Value(s string, opts Options) string {
	runes := []rune(s)
	if len(runes) <= opts.MaxLen {
		return s
	}
	truncated := string(runes[:opts.MaxLen]) + opts.Suffix
	if opts.ShowLength {
		truncated += fmt.Sprintf(" [%d]", len(runes))
	}
	return truncated
}

// Map applies Value to every entry in env, returning a new map with truncated
// values. The original map is not modified.
func Map(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Value(v, opts)
	}
	return out
}
