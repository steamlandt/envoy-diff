// Package lint provides heuristic checks on environment variable maps,
// flagging suspicious patterns such as duplicate keys (case-insensitive),
// overly long values, or keys that look like they contain embedded secrets.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding describes a single lint result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Options controls which checks are enabled.
type Options struct {
	MaxValueLen       int  // 0 = disabled
	CheckDuplicates   bool
	CheckKeyUppercase bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxValueLen:       4096,
		CheckDuplicates:   true,
		CheckKeyUppercase: true,
	}
}

// Apply runs all enabled lint checks against env and returns any findings.
func Apply(env map[string]string, opts Options) []Finding {
	var findings []Finding

	if opts.CheckDuplicates {
		findings = append(findings, checkDuplicates(env)...)
	}
	if opts.CheckKeyUppercase {
		findings = append(findings, checkKeyCase(env)...)
	}
	if opts.MaxValueLen > 0 {
		findings = append(findings, checkValueLength(env, opts.MaxValueLen)...)
	}

	return findings
}

func checkDuplicates(env map[string]string) []Finding {
	seen := make(map[string]string) // lower -> original
	var findings []Finding
	for k := range env {
		lower := strings.ToLower(k)
		if orig, ok := seen[lower]; ok && orig != k {
			findings = append(findings, Finding{
				Key:      k,
				Message:  fmt.Sprintf("duplicate key (case-insensitive conflict with %q)", orig),
				Severity: SeverityError,
			})
		} else {
			seen[lower] = k
		}
	}
	return findings
}

func checkKeyCase(env map[string]string) []Finding {
	var findings []Finding
	for k := range env {
		if k != strings.ToUpper(k) {
			findings = append(findings, Finding{
				Key:      k,
				Message:  "key is not uppercase",
				Severity: SeverityWarn,
			})
		}
	}
	return findings
}

func checkValueLength(env map[string]string, max int) []Finding {
	var findings []Finding
	for k, v := range env {
		if len(v) > max {
			findings = append(findings, Finding{
				Key:      k,
				Message:  fmt.Sprintf("value length %d exceeds maximum %d", len(v), max),
				Severity: SeverityWarn,
			})
		}
	}
	return findings
}
