// Package resolve provides variable interpolation for environment maps,
// expanding references like ${VAR} or $VAR within values.
package resolve

import (
	"os"
	"strings"
)

// Options controls how interpolation is performed.
type Options struct {
	// Fallback to the current process environment when a key is not found
	// in the provided map.
	FallbackToOS bool
	// MaxDepth limits recursive expansion to prevent infinite loops.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS: true,
		MaxDepth:     10,
	}
}

// Apply expands variable references in every value of env using the keys
// defined within env itself (and optionally the OS environment).
// It returns a new map and does not mutate the input.
func Apply(env map[string]string, opts Options) map[string]string {
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = 10
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = expand(v, env, opts, 0)
	}
	return out
}

// expand recursively resolves variable references in s.
func expand(s string, env map[string]string, opts Options, depth int) string {
	if depth >= opts.MaxDepth {
		return s
	}
	return os.Expand(s, func(key string) string {
		if val, ok := env[key]; ok {
			return expand(val, env, opts, depth+1)
		}
		if opts.FallbackToOS {
			return os.Getenv(key)
		}
		return ""
	})
}

// ApplyToValue expands a single value string using the provided env map.
func ApplyToValue(value string, env map[string]string, opts Options) string {
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = 10
	}
	return expand(value, env, opts, 0)
}

// HasReferences reports whether s contains any $VAR or ${VAR} references.
func HasReferences(s string) bool {
	return strings.ContainsRune(s, '$')
}
