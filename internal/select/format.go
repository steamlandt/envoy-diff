package sel

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

const (
	ansiGreen = "\033[32m"
	ansiReset = "\033[0m"
	ansiGray  = "\033[90m"
)

// FormatText writes a human-readable summary of the selection result to w.
// color enables ANSI colour codes.
func FormatText(w io.Writer, selected []diff.Entry, total int, color bool) {
	dropped := total - len(selected)

	if len(selected) == 0 {
		fmt.Fprintln(w, paint("no entries selected", ansiGray, color))
		return
	}

	for _, e := range selected {
		fmt.Fprintf(w, "  %s\n", paint(e.Key, ansiGreen, color))
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n",
		paint(fmt.Sprintf("%d selected, %d dropped", len(selected), dropped), ansiGray, color),
	)
}

// ParseKeys splits a comma-separated key list into individual key strings,
// trimming whitespace from each entry.
func ParseKeys(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func paint(s, code string, color bool) string {
	if !color {
		return s
	}
	return code + s + ansiReset
}
