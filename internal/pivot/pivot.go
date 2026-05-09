// Package pivot provides utilities for transposing environment variable
// diff results into a key-centric view, grouping entries by key across
// multiple named sources.
package pivot

import (
	"sort"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// SourceEntry holds the value of a key as seen in one named source.
type SourceEntry struct {
	Source string
	Value  string
	Present bool
}

// Row represents a single key across all sources.
type Row struct {
	Key     string
	Sources []SourceEntry
	Uniform bool // true when every present source agrees on the value
}

// Table is the full pivot result.
type Table struct {
	Sources []string
	Rows    []Row
}

// Build constructs a Table from a map of named env maps.
// sources defines the display order of source names.
func Build(sources []string, envs map[string]map[string]string) Table {
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]Row, 0, len(keys))
	for _, key := range keys {
		entries := make([]SourceEntry, len(sources))
		var firstVal string
		presentCount := 0
		uniform := true

		for i, src := range sources {
			val, ok := envs[src][key]
			entries[i] = SourceEntry{Source: src, Value: val, Present: ok}
			if ok {
				if presentCount == 0 {
					firstVal = val
				} else if val != firstVal {
					uniform = false
				}
				presentCount++
			}
		}
		if presentCount < 2 {
			uniform = false
		}
		rows = append(rows, Row{Key: key, Sources: entries, Uniform: uniform})
	}

	return Table{Sources: sources, Rows: rows}
}

// FromDiffs builds a Table from a slice of (sourceName, []diff.Entry) pairs.
// It reconstructs per-source maps from the diff entries relative to an empty base.
func FromDiffs(sources []string, diffs map[string][]diff.Entry) Table {
	envs := make(map[string]map[string]string, len(sources))
	for _, src := range sources {
		m := map[string]string{}
		for _, e := range diffs[src] {
			if e.Status != diff.Removed {
				m[e.Key] = e.New
			}
		}
		envs[src] = m
	}
	return Build(sources, envs)
}
