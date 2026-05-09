// Package select provides key-based selection and exclusion of environment entries.
package sel

import (
	"regexp"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Options controls which entries are kept or dropped.
type Options struct {
	// Keys is an explicit list of keys to keep. Empty means keep all.
	Keys []string
	// Exclude is an explicit list of keys to drop.
	Exclude []string
	// Pattern is a regex applied to key names; only matching keys are kept.
	Pattern string
	// ExcludePattern is a regex applied to key names; matching keys are dropped.
	ExcludePattern string
}

// Apply returns a filtered copy of entries according to opts.
func Apply(entries []diff.Entry, opts Options) ([]diff.Entry, error) {
	var includeRe, excludeRe *regexp.Regexp
	var err error

	if opts.Pattern != "" {
		includeRe, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}
	if opts.ExcludePattern != "" {
		excludeRe, err = regexp.Compile(opts.ExcludePattern)
		if err != nil {
			return nil, err
		}
	}

	keySet := toSet(opts.Keys)
	excludeSet := toSet(opts.Exclude)

	result := make([]diff.Entry, 0, len(entries))
	for _, e := range entries {
		if !keep(e.Key, keySet, excludeSet, includeRe, excludeRe) {
			continue
		}
		result = append(result, e)
	}
	return result, nil
}

func keep(key string, include, exclude map[string]struct{}, includeRe, excludeRe *regexp.Regexp) bool {
	if len(exclude) > 0 {
		if _, ok := exclude[strings.ToUpper(key)]; ok {
			return false
		}
	}
	if excludeRe != nil && excludeRe.MatchString(key) {
		return false
	}
	if len(include) > 0 {
		if _, ok := include[strings.ToUpper(key)]; !ok {
			return false
		}
	}
	if includeRe != nil && !includeRe.MatchString(key) {
		return false
	}
	return true
}

func toSet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToUpper(k)] = struct{}{}
	}
	return m
}
