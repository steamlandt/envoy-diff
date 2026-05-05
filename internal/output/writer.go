// Package output handles writing diff results to various destinations.
package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envoy-diff/internal/diff"
)

// Format represents the output format for diff results.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options configures output behaviour.
type Options struct {
	Format  Format
	NoColor bool
	Out     io.Writer
}

// DefaultOptions returns Options writing plain text to stdout.
func DefaultOptions() Options {
	return Options{
		Format:  FormatText,
		NoColor: false,
		Out:     os.Stdout,
	}
}

// Write renders the diff result according to opts and writes it to opts.Out.
// It returns the number of differing entries and any write error.
func Write(result []diff.Entry, opts Options) (int, error) {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	switch opts.Format {
	case FormatJSON:
		return writeJSON(result, opts.Out)
	default:
		return writeText(result, opts)
	}
}

func writeText(result []diff.Entry, opts Options) (int, error) {
	formatted := diff.Format(result, !opts.NoColor)
	if formatted == "" {
		_, err := fmt.Fprintln(opts.Out, "No differences found.")
		return 0, err
	}
	_, err := fmt.Fprint(opts.Out, formatted)
	return countChanged(result), err
}

func countChanged(result []diff.Entry) int {
	n := 0
	for _, e := range result {
		if e.Status != diff.StatusEqual {
			n++
		}
	}
	return n
}
