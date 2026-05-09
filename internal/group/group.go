// Package group provides utilities for grouping environment variable entries
// by a key prefix, enabling structured views of namespaced configuration.
package group

import (
	"sort"
	"strings"
)

// Entry represents a single environment variable key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Group holds entries that share a common prefix.
type Group struct {
	Prefix  string
	Entries []Entry
}

// Options controls how grouping is performed.
type Options struct {
	// Separator is the delimiter used to split key prefixes (default "_").
	Separator string
	// MinGroupSize skips groups with fewer entries than this value.
	MinGroupSize int
}

// DefaultOptions returns sensible defaults for grouping.
func DefaultOptions() Options {
	return Options{
		Separator:    "_",
		MinGroupSize: 1,
	}
}

// Apply groups the provided map of environment variables by their first
// prefix segment (split on opts.Separator). Keys without a separator are
// placed in a group named "".
func Apply(env map[string]string, opts Options) []Group {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	buckets := make(map[string][]Entry)
	for k, v := range env {
		prefix := ""
		if idx := strings.Index(k, opts.Separator); idx > 0 {
			prefix = k[:idx]
		}
		buckets[prefix] = append(buckets[prefix], Entry{Key: k, Value: v})
	}

	var groups []Group
	for prefix, entries := range buckets {
		if len(entries) < opts.MinGroupSize {
			continue
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
		groups = append(groups, Group{Prefix: prefix, Entries: entries})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})
	return groups
}
