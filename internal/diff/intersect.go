package diff

// IntersectOptions controls the behaviour of Intersect.
type IntersectOptions struct {
	// OnlyChanged, when true, returns only keys whose values differ between A
	// and B.  When false every key present in both maps is returned.
	OnlyChanged bool
}

// DefaultIntersectOptions returns conservative defaults.
func DefaultIntersectOptions() IntersectOptions {
	return IntersectOptions{
		OnlyChanged: false,
	}
}

// IntersectResult holds a single key that appears in both input maps.
type IntersectResult struct {
	Key     string
	ValueA  string
	ValueB  string
	Changed bool
}

// Intersect returns all keys that are present in both a and b, together with
// their respective values and a flag indicating whether the values differ.
// The result slice is ordered by key.
func Intersect(a, b map[string]string, opts IntersectOptions) []IntersectResult {
	var results []IntersectResult

	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			continue
		}
		changed := va != vb
		if opts.OnlyChanged && !changed {
			continue
		}
		results = append(results, IntersectResult{
			Key:     k,
			ValueA:  va,
			ValueB:  vb,
			Changed: changed,
		})
	}

	// stable, deterministic order
	sortIntersect(results)
	return results
}

func sortIntersect(rs []IntersectResult) {
	// insertion sort — slices are small in practice
	for i := 1; i < len(rs); i++ {
		for j := i; j > 0 && rs[j].Key < rs[j-1].Key; j-- {
			rs[j], rs[j-1] = rs[j-1], rs[j]
		}
	}
}
