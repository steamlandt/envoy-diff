// Package merge provides utilities for merging multiple environment variable
// maps into a single map, with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first source that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last source that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears in multiple sources.
	StrategyError
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	Values []string // one entry per source index that defined the key
}

// Result holds the merged environment map and any conflicts that were detected.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines the provided environment maps according to the given strategy.
// Sources are processed in order; index 0 is considered the "first" source.
func Merge(strategy Strategy, sources ...map[string]string) (Result, error) {
	env := make(map[string]string)
	// track which source indices contributed each key
	origins := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			origins[k] = append(origins[k], v)
		}
	}

	var conflicts []Conflict

	for k, vals := range origins {
		if len(vals) == 1 {
			env[k] = vals[0]
			continue
		}
		// conflict
		switch strategy {
		case StrategyFirst:
			env[k] = vals[0]
		case StrategyLast:
			env[k] = vals[len(vals)-1]
		case StrategyError:
			return Result{}, fmt.Errorf("merge: key %q defined in multiple sources", k)
		}
		conflicts = append(conflicts, Conflict{Key: k, Values: vals})
	}

	return Result{Env: env, Conflicts: conflicts}, nil
}
