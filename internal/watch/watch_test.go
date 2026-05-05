package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envoy-diff/internal/watch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := watch.New("/nonexistent/path/.env", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\n")

	w, err := watch.New(p, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Stop()

	// Modify the file after a short delay.
	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(p, []byte("FOO=changed\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-w.C:
		if ev.Path != p {
			t.Errorf("event path = %q, want %q", ev.Path, p)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected OldHash != NewHash after change")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoSpuriousEvents(t *testing.T) {
	p := writeTempEnv(t, "FOO=stable\n")

	w, err := watch.New(p, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Stop()

	// Wait without modifying the file.
	time.Sleep(120 * time.Millisecond)

	select {
	case ev := <-w.C:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	default:
		// expected: no event
	}
}

func TestWatcher_Stop(t *testing.T) {
	p := writeTempEnv(t, "A=1\n")

	w, err := watch.New(p, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	w.Stop()

	// After stop, modifying the file should not deliver events.
	time.Sleep(40 * time.Millisecond)
	_ = os.WriteFile(p, []byte("A=2\n"), 0o600)
	time.Sleep(80 * time.Millisecond)

	select {
	case ev := <-w.C:
		t.Errorf("received event after Stop: %+v", ev)
	default:
	}
}
