// Package rename provides utilities for renaming environment variable keys
// across a map, supporting bulk renames via a mapping table.
package rename

import (
	"errors"
	"fmt"
)

// ErrDuplicateTarget is returned when two source keys map to the same target key.
var ErrDuplicateTarget = errors.New("rename: duplicate target key")

// Options controls the behaviour of Apply.
type Options struct {
	// FailOnMissing causes Apply to return an error if a source key in the
	// mapping does not exist in the input map.
	FailOnMissing bool
	// Overwrite allows the renamed key to overwrite an existing key in the
	// output. When false, Apply returns an error if the target already exists.
	Overwrite bool
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
		Overwrite:     false,
	}
}

// Apply returns a new map with keys renamed according to mapping.
// The original map is never mutated.
func Apply(env map[string]string, mapping map[string]string, opts Options) (map[string]string, error) {
	// Validate that no two sources map to the same target.
	seen := make(map[string]string, len(mapping))
	for src, dst := range mapping {
		if prev, ok := seen[dst]; ok {
			return nil, fmt.Errorf("%w: %q and %q both map to %q", ErrDuplicateTarget, prev, src, dst)
		}
		seen[dst] = src
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for src, dst := range mapping {
		val, exists := out[src]
		if !exists {
			if opts.FailOnMissing {
				return nil, fmt.Errorf("rename: source key %q not found", src)
			}
			continue
		}
		if _, targetExists := out[dst]; targetExists && !opts.Overwrite {
			return nil, fmt.Errorf("rename: target key %q already exists", dst)
		}
		delete(out, src)
		out[dst] = val
	}
	return out, nil
}
