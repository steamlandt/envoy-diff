// Package transform applies key/value transformations to environment variable maps.
package transform

import (
	"fmt"
	"strings"
)

// Options controls how transformations are applied.
type Options struct {
	// KeyPrefix prepends a string to every key.
	KeyPrefix string
	// KeySuffix appends a string to every key.
	KeySuffix string
	// UpperKeys converts all keys to uppercase.
	UpperKeys bool
	// LowerKeys converts all keys to lowercase.
	LowerKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
}

// DefaultOptions returns an Options with no transformations enabled.
func DefaultOptions() Options {
	return Options{}
}

// Result holds the outcome of a single key transformation.
type Result struct {
	OriginalKey string
	NewKey      string
	Value       string
	Skipped     bool
	Reason      string
}

// Apply transforms the given env map according to opts.
// It returns a new map and a slice of Result describing each transformation.
func Apply(env map[string]string, opts Options) (map[string]string, []Result, error) {
	if opts.UpperKeys && opts.LowerKeys {
		return nil, nil, fmt.Errorf("transform: UpperKeys and LowerKeys are mutually exclusive")
	}

	out := make(map[string]string, len(env))
	results := make([]Result, 0, len(env))

	for k, v := range env {
		newKey := k
		if opts.UpperKeys {
			newKey = strings.ToUpper(newKey)
		} else if opts.LowerKeys {
			newKey = strings.ToLower(newKey)
		}
		if opts.KeyPrefix != "" {
			newKey = opts.KeyPrefix + newKey
		}
		if opts.KeySuffix != "" {
			newKey = newKey + opts.KeySuffix
		}

		newVal := v
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		if _, exists := out[newKey]; exists {
			results = append(results, Result{
				OriginalKey: k,
				NewKey:      newKey,
				Value:       newVal,
				Skipped:     true,
				Reason:      "key collision after transformation",
			})
			continue
		}

		out[newKey] = newVal
		results = append(results, Result{
			OriginalKey: k,
			NewKey:      newKey,
			Value:       newVal,
		})
	}
	return out, results, nil
}
