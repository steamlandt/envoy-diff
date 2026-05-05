// Package export provides functionality to serialize environment variable
// diff results into shell-sourceable or dotenv-compatible output formats.
package export

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yourorg/envoy-diff/internal/diff"
)

// Format represents the output format for exported variables.
type Format string

const (
	// FormatShell produces POSIX shell export statements.
	FormatShell Format = "shell"
	// FormatDotenv produces KEY=VALUE dotenv-compatible lines.
	FormatDotenv Format = "dotenv"
)

// Options controls the behaviour of Write.
type Options struct {
	// Format selects the output format. Defaults to FormatShell.
	Format Format
	// OnlyChanged limits output to keys that changed or were added.
	OnlyChanged bool
}

// Write serialises the environment map derived from diff entries to w.
// When opts.OnlyChanged is true only Added/Changed entries are emitted.
func Write(w io.Writer, entries []diff.Entry, opts Options) error {
	if opts.Format == "" {
		opts.Format = FormatShell
	}

	keys := make([]string, 0, len(entries))
	index := make(map[string]diff.Entry, len(entries))
	for _, e := range entries {
		if opts.OnlyChanged && e.Status == diff.StatusUnchanged {
			continue
		}
		if e.Status == diff.StatusRemoved {
			continue
		}
		keys = append(keys, e.Key)
		index[e.Key] = e
	}
	sort.Strings(keys)

	for _, k := range keys {
		e := index[k]
		val := quote(e.ValueB)
		var line string
		switch opts.Format {
		case FormatDotenv:
			line = fmt.Sprintf("%s=%s\n", k, val)
		default:
			line = fmt.Sprintf("export %s=%s\n", k, val)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

// quote wraps v in single quotes if it contains shell-special characters.
func quote(v string) string {
	if v == "" {
		return "\"\""
	}
	if strings.ContainsAny(v, " \t\n$`\\\"|&;()<>!#") {
		return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
	}
	return v
}
