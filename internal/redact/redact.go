// Package redact provides utilities for masking sensitive environment
// variable values before display or output.
package redact

import "strings"

// DefaultPatterns is the list of key substrings that trigger redaction
// when no custom patterns are supplied.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
}

const mask = "***REDACTED***"

// Apply returns a copy of env with sensitive values replaced by the mask
// string. Keys are matched case-insensitively against patterns. If patterns
// is nil or empty, DefaultPatterns is used.
func Apply(env map[string]string, patterns []string) map[string]string {
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive reports whether key contains any of the given patterns
// (case-insensitive substring match).
func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Mask returns the sentinel string used to replace redacted values.
func Mask() string { return mask }
