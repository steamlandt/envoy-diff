package validate

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// FormatText writes a human-readable report of issues to w.
// If color is true, ANSI escape codes are used.
func FormatText(w io.Writer, issues []Issue, color bool) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No validation issues found.")
		return
	}

	// Sort for deterministic output: errors first, then by key.
	sorted := make([]Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			return sorted[i].Severity == "error"
		}
		return sorted[i].Key < sorted[j].Key
	})

	for _, issue := range sorted {
		line := formatIssue(issue, color)
		fmt.Fprintln(w, line)
	}

	errors, warnings := countIssues(issues)
	fmt.Fprintf(w, "\n%d error(s), %d warning(s)\n", errors, warnings)
}

func formatIssue(i Issue, color bool) string {
	var prefix string
	switch i.Severity {
	case "error":
		prefix = applyColor("ERROR", colorRed, color)
	case "warning":
		prefix = applyColor("WARN ", colorYellow, color)
	default:
		prefix = strings.ToUpper(i.Severity)
	}
	return fmt.Sprintf("  %s  %-30s %s", prefix, i.Key, i.Message)
}

func applyColor(s, code string, enabled bool) string {
	if !enabled {
		return s
	}
	return code + s + colorReset
}

func countIssues(issues []Issue) (errors, warnings int) {
	for _, i := range issues {
		switch i.Severity {
		case "error":
			errors++
		case "warning":
			warnings++
		}
	}
	return
}
