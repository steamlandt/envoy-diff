package diff

// Result holds the outcome of comparing two env maps.
type Result struct {
	OnlyInA  map[string]string // keys present in A but not B
	OnlyInB  map[string]string // keys present in B but not A
	Changed  map[string][2]string // keys in both but with different values [A, B]
	Unchanged map[string]string // keys with identical values in both
}

// Compare returns a Result describing the differences between env maps a and b.
func Compare(a, b map[string]string) Result {
	r := Result{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				r.Unchanged[k] = va
			} else {
				r.Changed[k] = [2]string{va, vb}
			}
		} else {
				r.OnlyInA[k] = va
			}
	}

	for k, vb := range b {
		if _, ok := a[k]; !ok {
			r.OnlyInB[k] = vb
		}
	}

	return r
}

// HasDifferences returns true when there is at least one addition, removal, or change.
func (r Result) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Changed) > 0
}
