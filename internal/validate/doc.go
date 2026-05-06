// Package validate provides validation checks for environment variable maps.
//
// It detects common problems such as:
//   - Empty or blank keys
//   - Keys containing invalid characters (by default only [A-Za-z0-9_] are allowed)
//   - Values that exceed a configurable maximum length
//   - Empty values (reported as warnings, not errors)
//
// Usage:
//
//	opts := validate.DefaultOptions()
//	issues := validate.Apply(envMap, opts)
//	if validate.HasErrors(issues) {
//	    validate.FormatText(os.Stderr, issues, true)
//	    os.Exit(1)
//	}
//
// Validation is non-destructive: the input map is never modified.
package validate
