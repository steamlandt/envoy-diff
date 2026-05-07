// Package scope provides utilities for scoping environment variable maps
// to a named prefix, allowing logical namespacing of keys.
package scope

import (
	"strings"
)

// Options controls the behaviour of scope operations.
type Options struct {
	// Separator is placed between the prefix and the key name.
	// Defaults to "_".
	Separator string

	// StripOnExtract removes the prefix from keys when extracting.
	// Defaults to true.
	StripOnExtract bool
}

// DefaultOptions returns sensible defaults for scope operations.
func DefaultOptions() Options {
	return Options{
		Separator:      "_",
		StripOnExtract: true,
	}
}

// Inject prefixes every key in env with "<prefix><sep>" and returns
// the resulting map. The original map is not mutated.
func Inject(prefix string, env map[string]string, opts Options) map[string]string {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	out := make(map[string]string, len(env))
	pfx := strings.ToUpper(prefix) + opts.Separator
	for k, v := range env {
		out[pfx+k] = v
	}
	return out
}

// Extract returns only the entries whose keys begin with
// "<prefix><sep>". If StripOnExtract is true the prefix is removed
// from the returned keys. The original map is not mutated.
func Extract(prefix string, env map[string]string, opts Options) map[string]string {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	out := make(map[string]string)
	pfx := strings.ToUpper(prefix) + opts.Separator
	for k, v := range env {
		if strings.HasPrefix(k, pfx) {
			key := k
			if opts.StripOnExtract {
				key = k[len(pfx):]
			}
			if key != "" {
				out[key] = v
			}
		}
	}
	return out
}
