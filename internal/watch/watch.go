// Package watch provides functionality to monitor .env files for changes
// and emit notifications when the environment variable set is modified.
package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a file change notification.
type Event struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a file for changes and sends events on the C channel.
type Watcher struct {
	C        <-chan Event
	Done     chan struct{}
	path     string
	interval time.Duration
	lastHash string
}

// New creates a new Watcher for the given file path, polling at the given interval.
// The watcher begins polling immediately. Call Stop to shut it down.
func New(path string, interval time.Duration) (*Watcher, error) {
	h, err := hashFile(path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial hash: %w", err)
	}

	ch := make(chan Event, 1)
	w := &Watcher{
		C:        ch,
		Done:     make(chan struct{}),
		path:     path,
		interval: interval,
		lastHash: h,
	}

	go w.poll(ch)
	return w, nil
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.Done)
}

func (w *Watcher) poll(ch chan<- Event) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.Done:
			return
		case <-ticker.C:
			h, err := hashFile(w.path)
			if err != nil || h == w.lastHash {
				continue
			}
			old := w.lastHash
			w.lastHash = h
			select {
			case ch <- Event{Path: w.path, OldHash: old, NewHash: h}:
			default:
			}
		}
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
