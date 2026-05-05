package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	sourceA  string
	sourceB  string
	color    bool
	onlyKeys bool
	exitCode bool
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("envoy-diff", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &config{}

	fs.BoolVar(&cfg.color, "color", isTerminal(), "colorize output")
	fs.BoolVar(&cfg.onlyKeys, "keys", false, "show only key names, not values")
	fs.BoolVar(&cfg.exitCode, "exit-code", false, "exit with code 2 if differences found")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("%w\n%s", err, usage())
	}

	if fs.NArg() != 2 {
		return nil, errors.New(usage())
	}

	cfg.sourceA = fs.Arg(0)
	cfg.sourceB = fs.Arg(1)

	return cfg, nil
}

func usage() string {
	return `Usage: envoy-diff [flags] <source-a> <source-b>

Sources can be:
  path/to/file.env   — a .env file
  self               — the current process environment
  pid:<N>            — environment of process with PID N

Flags:
  -color        colorize output (default: auto)
  -keys         show only key names, omit values
  -exit-code    exit with code 2 when differences are found`
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
