// Package mask provides utilities for partially obscuring environment
// variable values in output, preserving a configurable number of leading
// characters so that values remain identifiable without being fully exposed.
package mask

import "strings"

// DefaultVisibleChars is the number of characters shown before masking.
const DefaultVisibleChars = 4

// DefaultMaskChar is the character used to replace hidden characters.
const DefaultMaskChar = '*'

// Options controls how masking is applied.
type Options struct {
	// VisibleChars is the number of leading characters to reveal.
	VisibleChars int
	// MaskChar is the character used to obscure the rest of the value.
	MaskChar rune
	// MinLength is the minimum value length before masking is applied.
	// Values shorter than this are masked entirely.
	MinLength int
}

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		VisibleChars: DefaultVisibleChars,
		MaskChar:     DefaultMaskChar,
		MinLength:    8,
	}
}

// Value masks a single string value according to opts.
func Value(v string, opts Options) string {
	if v == "" {
		return ""
	}
	if len(v) < opts.MinLength {
		return strings.Repeat(string(opts.MaskChar), len(v))
	}
	visible := opts.VisibleChars
	if visible >= len(v) {
		visible = len(v) - 1
	}
	if visible < 0 {
		visible = 0
	}
	return v[:visible] + strings.Repeat(string(opts.MaskChar), len(v)-visible)
}

// Map applies masking to every value in the provided map, returning a new map.
// The original map is not modified.
func Map(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Value(v, opts)
	}
	return out
}
