// Package source provides a unified abstraction for loading environment
// variables from heterogeneous origins.
//
// Supported source kinds:
//
//	- File  — parses a .env file on disk via the parser package.
//	- PID   — reads /proc/<pid>/environ (Linux) via the process package.
//	- Self  — reads the current process's own environment.
//
// Usage:
//
//	env, err := source.Load(source.Source{Kind: source.File, Path: ".env"})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The Load function returns a plain map[string]string that can be passed
// directly to diff.Compare.
package source
