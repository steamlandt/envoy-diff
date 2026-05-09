// Package pivot transposes environment variable data from a
// source-centric view ("what changed in source X?") into a
// key-centric view ("how does KEY differ across all sources?").
//
// # Overview
//
// Given N named sources (files, PIDs, snapshots), Build constructs a
// Table where each Row corresponds to one environment variable key and
// contains the value—or absence marker—for every source.
//
// Rows are sorted alphabetically by key. The Uniform field on each Row
// is true only when every source that contains the key agrees on the
// same value, making it easy to filter to only the interesting rows.
//
// # Usage
//
//	sources := []string{".env.staging", ".env.production"}
//	envs := map[string]map[string]string{
//		".env.staging":    stagingMap,
//		".env.production": prodMap,
//	}
//	table := pivot.Build(sources, envs)
//	pivot.FormatText(os.Stdout, table, pivot.DefaultFormatOptions())
package pivot
