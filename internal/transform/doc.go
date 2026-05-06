// Package transform provides utilities for renaming and normalising keys
// and values in an environment variable map.
//
// Supported transformations:
//   - Prefix / suffix injection on keys
//   - Uppercase or lowercase key normalisation
//   - Whitespace trimming on values
//
// Collisions that arise after key transformations are reported as skipped
// entries rather than silently overwriting existing keys.
package transform
