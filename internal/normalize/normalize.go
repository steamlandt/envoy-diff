// Package normalize provides utilities for normalizing environment variable
// maps — trimming whitespace from keys and values, deduplicating entries,
// and optionally uppercasing all keys.
package normalize

import (
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from keys and values.
	TrimSpace bool

	// UpperKeys converts all keys to uppercase.
	UpperKeys bool

	// DeduplicateKeys keeps only the last occurrence of a duplicate key.
	// When false, duplicates are kept in insertion order.
	DeduplicateKeys bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:       true,
		UpperKeys:       false,
		DeduplicateKeys: true,
	}
}

// Result holds the output of a normalization pass.
type Result struct {
	// Env is the normalized map.
	Env map[string]string

	// Renamed lists keys that were renamed due to uppercasing.
	Renamed []string

	// Dropped lists keys that were removed as duplicates.
	Dropped []string
}

// Apply normalizes env according to opts and returns a Result.
func Apply(env map[string]string, opts Options) Result {
	out := make(map[string]string, len(env))
	seen := make(map[string]string, len(env)) // canonical key -> original key
	var renamed, dropped []string

	for k, v := range env {
		nk := k
		nv := v

		if opts.TrimSpace {
			nk = strings.TrimSpace(nk)
			nv = strings.TrimSpace(nv)
		}

		if opts.UpperKeys {
			upper := strings.ToUpper(nk)
			if upper != nk {
				renamed = append(renamed, nk)
			}
			nk = upper
		}

		if opts.DeduplicateKeys {
			if prev, exists := seen[nk]; exists {
				dropped = append(dropped, prev)
			}
			seen[nk] = nk
		}

		out[nk] = nv
	}

	return Result{
		Env:     out,
		Renamed: renamed,
		Dropped: dropped,
	}
}
