// Package lint implements heuristic checks for environment variable maps.
//
// It detects common mistakes such as:
//   - Keys that differ only in case (likely duplicates)
//   - Lowercase or mixed-case keys (unconventional for env vars)
//   - Values that exceed a configurable maximum length
//
// Usage:
//
//	findings := lint.Apply(env, lint.DefaultOptions())
//	lint.FormatText(os.Stderr, findings, isTerminal(os.Stderr))
//	if lint.HasErrors(findings) {
//		os.Exit(1)
//	}
package lint
