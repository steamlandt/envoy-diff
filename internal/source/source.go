// Package source resolves environment variable maps from various sources:
// .env files, running process PIDs, the current process, or saved snapshots.
package source

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/envoy-diff/internal/parser"
	"github.com/user/envoy-diff/internal/process"
	"github.com/user/envoy-diff/internal/snapshot"
)

// Load resolves an environment map from the given source string.
//
// Source formats:
//   - ""         — current process environment
//   - "self"     — current process environment (explicit)
//   - "pid:<n>"  — environment of process with PID n
//   - "snap:<p>" — load a saved snapshot from file path p
//   - anything else is treated as a .env file path
func Load(src string) (map[string]string, error) {
	switch {
	case src == "" || src == "self":
		return process.ReadSelf(), nil

	case strings.HasPrefix(src, "pid:"):
		pidStr := strings.TrimPrefix(src, "pid:")
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return nil, fmt.Errorf("source: invalid pid %q: %w", pidStr, err)
		}
		return process.ReadPID(pid)

	case strings.HasPrefix(src, "snap:"):
		path := strings.TrimPrefix(src, "snap:")
		snap, err := snapshot.Load(path)
		if err != nil {
			return nil, fmt.Errorf("source: load snapshot: %w", err)
		}
		return snap.Env, nil

	default:
		return parser.ParseFile(src)
	}
}
