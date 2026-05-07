// Package diff compares two environment variable maps and produces a list of
// Entry values describing what was added, removed, changed, or left unchanged.
//
// # Comparing
//
// Use [Compare] to produce a []Entry from two map[string]string values.
//
// # Formatting
//
// [Format] renders a []Entry as a human-readable, optionally colourised string.
//
// # Patching
//
// [Patch] serialises the changed entries into a unified, shell, or dotenv
// patch string that can be shared or stored.
//
// # Applying
//
// [Apply] takes a base environment map and a []Entry and returns a new map
// that reflects the target ("B") side of the diff, with optional conflict
// detection and skip-remove semantics.
//
// # Filtering
//
// [Filter] accepts a []Entry and a predicate function, returning only the
// entries for which the predicate returns true. This is useful for narrowing
// results to a specific [Kind] (e.g. only additions or removals) before
// passing them to [Format], [Patch], or [Apply].
package diff
