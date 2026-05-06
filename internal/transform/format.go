package transform

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatText writes a human-readable summary of transform results to w.
// If color is true, ANSI codes are included.
func FormatText(w io.Writer, results []Result, color bool) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No transformations applied.")
		return
	}

	// stable output
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].OriginalKey < sorted[j].OriginalKey
	})

	changed, skipped := 0, 0
	for _, r := range sorted {
		if r.Skipped {
			skipped++
			label := applyColor("SKIP", colorGray, color)
			fmt.Fprintf(w, "  %s  %s → %s (%s)\n", label, r.OriginalKey, r.NewKey, r.Reason)
		} else if r.OriginalKey != r.NewKey {
			changed++
			label := applyColor("RENAME", colorYellow, color)
			fmt.Fprintf(w, "  %s  %s → %s\n", label, r.OriginalKey, r.NewKey)
		}
	}

	parts := []string{}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d renamed", changed))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", skipped))
	}
	if len(parts) == 0 {
		fmt.Fprintln(w, "No key changes (values may have been trimmed).")
	} else {
		fmt.Fprintf(w, "Summary: %s\n", strings.Join(parts, ", "))
	}
}

func applyColor(s, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return code + s + colorReset
}
