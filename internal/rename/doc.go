// Package rename provides key-renaming utilities for environment variable maps.
//
// The primary entry point is Apply, which accepts an env map and a mapping of
// old key names to new key names, returning a new map with the renames applied.
// The original map is never mutated.
//
// Example usage:
//
//	mapping := map[string]string{
//		"LEGACY_DB_HOST": "DATABASE_HOST",
//		"LEGACY_DB_PORT": "DATABASE_PORT",
//	}
//	renamed, err := rename.Apply(env, mapping, rename.DefaultOptions())
//
// ParseMapping can be used to build a mapping from a slice of "OLD=NEW" strings,
// which is convenient when accepting rename pairs from CLI flags or config files.
package rename
