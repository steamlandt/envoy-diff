// Package filter provides utilities for filtering environment variable maps
// by key prefix, suffix, or regex pattern before diffing.
package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// Options holds the filtering configuration.
type Options struct {
	Prefix  string
	Suffix  string
	Pattern string
}

// Apply filters the given environment map according to opts, returning a new
// map containing only the matching keys. If no filter option is set the
// original map is returned unchanged.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	if opts.Prefix == "" && opts.Suffix == "" && opts.Pattern == "" {
		return env, nil
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if matchKey(k, opts.Prefix, opts.Suffix, re) {
			out[k] = v
		}
	}
	return out, nil
}

func matchKey(key, prefix, suffix string, re *regexp.Regexp) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if suffix != "" && !strings.HasSuffix(key, suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
