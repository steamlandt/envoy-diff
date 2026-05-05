package diff

import (
	"fmt"
	"sort"
	"strings"
)

// FormatOptions controls how diff output is rendered.
type FormatOptions struct {
	Color    bool
	OnlyKeys bool
}

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// Format renders a slice of Entry values as a human-readable diff string.
func Format(entries []Entry, opts FormatOptions) string {
	if len(entries) == 0 {
		return "(no differences)\n"
	}

	var sb strings.Builder
	for _, key := range sortedKeys(entries) {
		for _, e := range entries {
			if e.Key != key {
				continue
			}
			sb.WriteString(formatEntry(e, opts))
		}
	}
	return sb.String()
}

func formatEntry(e Entry, opts FormatOptions) string {
	switch e.Kind {
	case OnlyInA:
		line := formatLine("-", e.Key, e.ValueA, opts)
		return colorize(line, colorRed, opts.Color)
	case OnlyInB:
		line := formatLine("+", e.Key, e.ValueB, opts)
		return colorize(line, colorGreen, opts.Color)
	case Changed:
		var sb strings.Builder
		sb.WriteString(colorize(formatLine("-", e.Key, e.ValueA, opts), colorRed, opts.Color))
		sb.WriteString(colorize(formatLine("+", e.Key, e.ValueB, opts), colorYellow, opts.Color))
		return sb.String()
	}
	return ""
}

func formatLine(prefix, key, value string, opts FormatOptions) string {
	if opts.OnlyKeys {
		return fmt.Sprintf("%s %s\n", prefix, key)
	}
	return fmt.Sprintf("%s %s=%s\n", prefix, key, value)
}

func colorize(s, color string, enabled bool) string {
	if !enabled {
		return s
	}
	return color + s + colorReset
}

func sortedKeys(entries []Entry) []string {
	seen := make(map[string]struct{})
	for _, e := range entries {
		seen[e.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Summary returns a short human-readable summary of the diff, reporting the
// count of added, removed, and changed keys.
func Summary(entries []Entry) string {
	var added, removed, changed int
	for _, e := range entries {
		switch e.Kind {
		case OnlyInB:
			added++
		case OnlyInA:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d changed", added, removed, changed)
}
