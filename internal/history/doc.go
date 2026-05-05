// Package history provides persistent storage for diff run history.
//
// Each time envoy-diff is invoked, callers may record the resulting diff
// entries along with metadata (sources, timestamp) using [Record]. Stored
// entries can be retrieved with [List] and old entries removed with [Prune].
//
// Data is persisted as a JSON file (history.json) inside a user-supplied
// directory, typically under the OS cache or config directory:
//
//	os.UserCacheDir() + "/envoy-diff"
//
// Example:
//
//	err := history.Record(dir, history.Entry{
//		Timestamp: time.Now(),
//		SourceA:   "file:.env",
//		SourceB:   "pid:1234",
//		Results:   entries,
//	})
package history
