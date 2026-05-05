// Package filter provides key-based filtering for environment variable maps.
//
// Filtering can be applied by key prefix, key suffix, or a regular expression
// pattern. Multiple criteria are combined with AND semantics — a key must
// satisfy every non-empty criterion to be included in the output.
//
// Example usage:
//
//	out, err := filter.Apply(env, filter.Options{
//		Prefix:  "APP_",
//		Pattern: ".*_HOST$",
//	})
//
// If no options are set, Apply returns the original map without copying it.
package filter
