// Package watch implements lightweight file-change detection for envoy-diff.
//
// It polls a target .env file at a configurable interval, computes a SHA-256
// digest of the file contents on each tick, and emits an Event on the exported
// channel whenever the digest changes. This allows the CLI (or any consumer)
// to react to live edits without relying on OS-specific filesystem notification
// APIs such as inotify or FSEvents.
//
// Basic usage:
//
//	w, err := watch.New(".env", 2*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Stop()
//
//	for ev := range w.C {
//		fmt.Printf(".env changed: %s -> %s\n", ev.OldHash[:8], ev.NewHash[:8])
//	}
package watch
