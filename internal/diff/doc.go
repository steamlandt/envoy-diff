// Package diff provides utilities for comparing two sets of environment
// variables and formatting the resulting differences for human consumption.
//
// Basic usage:
//
//	a := map[string]string{"FOO": "bar", "OLD": "gone"}
//	b := map[string]string{"FOO": "changed", "NEW": "here"}
//
//	result := diff.Compare(a, b)
//	diff.Format(os.Stdout, result, ".env.dev", ".env.prod")
//
// The Result type exposes four categories:
//   - OnlyInA  – keys removed relative to B
//   - OnlyInB  – keys added relative to A
//   - Changed  – keys present in both with differing values
//   - Unchanged – keys present in both with identical values
package diff
