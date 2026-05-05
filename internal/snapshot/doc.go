// Package snapshot provides save and load functionality for capturing
// environment variable sets to disk.
//
// A snapshot records a labelled, timestamped copy of an environment map
// sourced from a .env file or running process. Snapshots are stored as
// indented JSON and can be loaded back as a source for diffing against
// a live environment or another snapshot.
//
// Typical usage:
//
//	// Save current process environment
//	snapshot.Save(".snapshots/baseline.json", "baseline", "self", env)
//
//	// Later, load and diff
//	snap, _ := snapshot.Load(".snapshots/baseline.json")
//	diff.Compare(snap.Env, liveEnv)
package snapshot
