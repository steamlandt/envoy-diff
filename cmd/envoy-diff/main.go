package main

import (
	"fmt"
	"os"

	"github.com/user/envoy-diff/internal/diff"
	"github.com/user/envoy-diff/internal/source"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	a, err := source.Load(cfg.sourceA)
	if err != nil {
		return fmt.Errorf("loading source A (%q): %w", cfg.sourceA, err)
	}

	b, err := source.Load(cfg.sourceB)
	if err != nil {
		return fmt.Errorf("loading source B (%q): %w", cfg.sourceB, err)
	}

	result := diff.Compare(a, b)

	output := diff.Format(result, diff.FormatOptions{
		Color:    cfg.color,
		OnlyKeys: cfg.onlyKeys,
	})

	fmt.Print(output)

	if cfg.exitCode && len(result) > 0 {
		os.Exit(2)
	}

	return nil
}
