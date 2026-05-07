package diff

import "errors"

// ApplyOptions controls the behaviour of Apply.
type ApplyOptions struct {
	// SkipRemoved, when true, does not delete keys that are marked as removed.
	SkipRemoved bool
	// FailOnConflict returns an error when a changed key already has a value
	// in base that differs from the patch's expected ValueA.
	FailOnConflict bool
}

// DefaultApplyOptions returns a sensible default configuration.
func DefaultApplyOptions() ApplyOptions {
	return ApplyOptions{}
}

// Apply takes a base environment map and a slice of diff entries and produces
// a new map that reflects the "B" side of the diff.
// Only StatusAdded, StatusRemoved, and StatusChanged entries are acted upon;
// StatusUnchanged entries are ignored.
func Apply(base map[string]string, entries []Entry, opts ApplyOptions) (map[string]string, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			out[e.Key] = e.ValueB

		case StatusRemoved:
			if !opts.SkipRemoved {
				delete(out, e.Key)
			}

		case StatusChanged:
			if opts.FailOnConflict {
				if cur, ok := out[e.Key]; ok && cur != e.ValueA {
					return nil, errors.New("conflict on key " + e.Key)
				}
			}
			out[e.Key] = e.ValueB
		}
	}
	return out, nil
}
