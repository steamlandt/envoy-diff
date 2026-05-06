// Package convert provides utilities for transforming environment variable
// maps between different representations used across envoy-diff.
package convert

import (
	"fmt"
	"sort"
	"strings"
)

// ToSlice converts a map of environment variables to a slice of KEY=VALUE strings,
// sorted alphabetically by key.
func ToSlice(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(env))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return out
}

// FromSlice parses a slice of KEY=VALUE strings into a map.
// Entries without an '=' separator are stored with an empty string value.
// Blank entries are skipped.
func FromSlice(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[p] = ""
			continue
		}
		out[p[:idx]] = p[idx+1:]
	}
	return out
}

// MergeOverride returns a new map containing all keys from base, with values
// overridden by any matching keys present in override.
func MergeOverride(base, override map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}

// Keys returns a sorted slice of all keys in the map.
func Keys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
