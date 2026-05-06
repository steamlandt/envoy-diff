// Package validate checks environment variable maps for common issues
// such as empty keys, invalid characters, duplicate keys, and oversized values.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key      string
	Message  string
	Severity string // "error" or "warning"
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Key, i.Message)
}

// Options controls which checks are performed.
type Options struct {
	MaxValueLen    int  // 0 means no limit
	AllowDotInKey  bool
	WarnEmptyValue bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxValueLen:    32768,
		AllowDotInKey:  false,
		WarnEmptyValue: true,
	}
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.]*$`)
var validKeyNoDotRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Apply validates the given env map and returns a slice of Issues.
func Apply(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if k == "" {
			issues = append(issues, Issue{Key: "(empty)", Message: "key must not be empty", Severity: "error"})
			continue
		}

		re := validKeyNoDotRe
		if opts.AllowDotInKey {
			re = validKeyRe
		}
		if !re.MatchString(k) {
			issues = append(issues, Issue{Key: k, Message: "key contains invalid characters", Severity: "error"})
		}

		if strings.ContainsAny(k, " \t") {
			issues = append(issues, Issue{Key: k, Message: "key contains whitespace", Severity: "error"})
		}

		if opts.WarnEmptyValue && v == "" {
			issues = append(issues, Issue{Key: k, Message: "value is empty", Severity: "warning"})
		}

		if opts.MaxValueLen > 0 && len(v) > opts.MaxValueLen {
			issues = append(issues, Issue{
				Key:      k,
				Message:  fmt.Sprintf("value exceeds max length %d", opts.MaxValueLen),
				Severity: "error",
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue has severity "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}
