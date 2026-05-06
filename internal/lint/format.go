package lint

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiReset  = "\033[0m"
)

// FormatText writes a human-readable lint report to w.
// color controls whether ANSI escape codes are emitted.
func FormatText(w io.Writer, findings []Finding, color bool) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "No lint issues found.")
		return
	}

	// Sort for deterministic output: errors first, then by key.
	sorted := make([]Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			return sorted[i].Severity == SeverityError
		}
		return sorted[i].Key < sorted[j].Key
	})

	var errs, warns int
	for _, f := range sorted {
		prefix := severityLabel(f.Severity, color)
		fmt.Fprintf(w, "  %s %s: %s\n", prefix, f.Key, f.Message)
		if f.Severity == SeverityError {
			errs++
		} else {
			warns++
		}
	}

	parts := []string{}
	if errs > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", errs))
	}
	if warns > 0 {
		parts = append(parts, fmt.Sprintf("%d warning(s)", warns))
	}
	fmt.Fprintf(w, "\n%s\n", strings.Join(parts, ", "))
}

// HasErrors returns true if any finding has error severity.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}

func severityLabel(s Severity, color bool) string {
	switch s {
	case SeverityError:
		if color {
			return ansiRed + "[error]" + ansiReset
		}
		return "[error]"
	default:
		if color {
			return ansiYellow + "[warn] " + ansiReset
		}
		return "[warn] "
	}
}
