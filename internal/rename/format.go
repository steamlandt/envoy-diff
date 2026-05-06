package rename

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatText writes a human-readable summary of the rename mapping to w.
// It lists each rename as "OLD_KEY -> NEW_KEY" and notes keys that were
// skipped because they were absent from env.
func FormatText(w io.Writer, mapping map[string]string, env map[string]string, color bool) {
	if len(mapping) == 0 {
		fmt.Fprintln(w, "no renames defined")
		return
	}

	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, src := range keys {
		dst := mapping[src]
		_, exists := env[src]
		if exists {
			line := fmt.Sprintf("  %s -> %s", src, dst)
			if color {
				line = "\033[33m" + line + "\033[0m"
			}
			fmt.Fprintln(w, line)
		} else {
			skipped := fmt.Sprintf("  %s (skipped, not found)", src)
			if color {
				skipped = "\033[2m" + skipped + "\033[0m"
			}
			fmt.Fprintln(w, skipped)
		}
	}
}

// ParseMapping parses a slice of "OLD=NEW" strings into a rename mapping.
// Lines that do not contain "=" are ignored.
func ParseMapping(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			continue
		}
		src := strings.TrimSpace(p[:idx])
		dst := strings.TrimSpace(p[idx+1:])
		if src != "" && dst != "" {
			m[src] = dst
		}
	}
	return m
}
