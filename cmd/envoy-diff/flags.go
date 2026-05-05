package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envoy-diff/internal/filter"
)

// config holds all parsed CLI configuration.
type config struct {
	SourceA string
	SourceB string
	NoColor bool
	OnlyChanged bool
	Filter  filter.Options
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("envoy-diff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = usage

	var cfg config
	fs.BoolVar(&cfg.NoColor, "no-color", false, "disable coloured output")
	fs.BoolVar(&cfg.OnlyChanged, "only-changed", false, "show only changed/added/removed keys")
	fs.StringVar(&cfg.Filter.Prefix, "prefix", "", "only compare keys with this prefix")
	fs.StringVar(&cfg.Filter.Suffix, "suffix", "", "only compare keys with this suffix")
	fs.StringVar(&cfg.Filter.Pattern, "pattern", "", "only compare keys matching this regex")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() < 2 {
		return nil, errors.New("two sources required: <a> <b>")
	}

	cfg.SourceA = fs.Arg(0)
	cfg.SourceB = fs.Arg(1)
	return &cfg, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: envoy-diff [options] <source-a> <source-b>

Sources:
  path/to/file.env   read from a .env file
  self               read from the current process
  <pid>              read from a running process by PID

Options:`)
	fmt.Fprintln(os.Stderr, "  -no-color       disable coloured output")
	fmt.Fprintln(os.Stderr, "  -only-changed   show only changed/added/removed keys")
	fmt.Fprintln(os.Stderr, "  -prefix string  only compare keys with this prefix")
	fmt.Fprintln(os.Stderr, "  -suffix string  only compare keys with this suffix")
	fmt.Fprintln(os.Stderr, "  -pattern string only compare keys matching this regex")
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
