package diff

import (
	"fmt"
	"io"
	"sort"
)

// Format writes a human-readable diff to w using the provided source labels.
func Format(w io.Writer, r Result, labelA, labelB string) {
	if !r.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(r.OnlyInA) > 0 {
		fmt.Fprintf(w, "--- only in %s ---\n", labelA)
		for _, k := range sortedKeys(r.OnlyInA) {
			fmt.Fprintf(w, "  - %s=%s\n", k, r.OnlyInA[k])
		}
	}

	if len(r.OnlyInB) > 0 {
		fmt.Fprintf(w, "+++ only in %s +++\n", labelB)
		for _, k := range sortedKeys(r.OnlyInB) {
			fmt.Fprintf(w, "  + %s=%s\n", k, r.OnlyInB[k])
		}
	}

	if len(r.Changed) > 0 {
		fmt.Fprintln(w, "~~~ changed ~~~")
		keys := make([]string, 0, len(r.Changed))
		for k := range r.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pair := r.Changed[k]
			fmt.Fprintf(w, "  ~ %s\n", k)
			fmt.Fprintf(w, "      %s: %s\n", labelA, pair[0])
			fmt.Fprintf(w, "      %s: %s\n", labelB, pair[1])
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
