package pivot

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatOptions controls rendering behaviour.
type FormatOptions struct {
	Color       bool
	OnlyDiffing bool // skip rows where all sources agree
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{Color: true}
}

// FormatText writes a human-readable pivot table to w.
func FormatText(w io.Writer, t Table, opts FormatOptions) {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, "(no keys)")
		return
	}

	// header
	header := fmt.Sprintf("%-30s", "KEY")
	for _, src := range t.Sources {
		header += fmt.Sprintf("  %-24s", truncate(src, 24))
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", 30+26*len(t.Sources)))

	for _, row := range t.Rows {
		if opts.OnlyDiffing && row.Uniform {
			continue
		}

		line := fmt.Sprintf("%-30s", truncate(row.Key, 30))
		for _, se := range row.Sources {
			var cell string
			if !se.Present {
				cell = gray("(absent)", opts.Color)
			} else if !row.Uniform {
				cell = yellow(truncate(se.Value, 24), opts.Color)
			} else {
				cell = truncate(se.Value, 24)
			}
			line += fmt.Sprintf("  %-24s", cell)
		}
		fmt.Fprintln(w, line)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func yellow(s string, color bool) string {
	if !color {
		return s
	}
	return colorYellow + s + colorReset
}

func gray(s string, color bool) string {
	if !color {
		return s
	}
	return colorGray + s + colorReset
}
