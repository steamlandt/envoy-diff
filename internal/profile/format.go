package profile

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	colorGreen = "\033[32m"
	colorGray  = "\033[90m"
	colorReset = "\033[0m"
)

// FormatList writes a human-readable list of profile names to w.
func FormatList(w io.Writer, names []string, color bool) {
	if len(names) == 0 {
		fmt.Fprintln(w, "no profiles saved")
		return
	}
	for _, n := range names {
		if color {
			fmt.Fprintf(w, "  %s%s%s\n", colorGreen, n, colorReset)
		} else {
			fmt.Fprintf(w, "  %s\n", n)
		}
	}
}

// FormatProfile writes the key/value pairs of a profile to w.
func FormatProfile(w io.Writer, p Profile, color bool) {
	keys := make([]string, 0, len(p.Env))
	for k := range p.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if color {
		fmt.Fprintf(w, "%s# profile: %s  saved: %s%s\n",
			colorGray, p.Name, p.SavedAt.Format("2006-01-02 15:04:05"), colorReset)
	} else {
		fmt.Fprintf(w, "# profile: %s  saved: %s\n",
			p.Name, p.SavedAt.Format("2006-01-02 15:04:05"))
	}

	for _, k := range keys {
		v := p.Env[k]
		if strings.ContainsAny(v, " \t") {
			v = "\"" + v + "\""
		}
		if color {
			fmt.Fprintf(w, "%s%s%s=%s\n", colorGreen, k, colorReset, v)
		} else {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
}
