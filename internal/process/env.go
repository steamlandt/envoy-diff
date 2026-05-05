// Package process provides utilities for reading environment variables
// from running processes via /proc or OS-native APIs.
package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadPID returns the environment variables of the process with the given PID
// as a map of key→value pairs.
//
// On Linux it reads /proc/<pid>/environ; on other platforms it returns an
// error indicating the feature is unsupported.
func ReadPID(pid int) (map[string]string, error) {
	path := fmt.Sprintf("/proc/%d/environ", pid)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("process %d not found (does /proc exist on this OS?)", pid)
		}
		return nil, fmt.Errorf("read environ for pid %d: %w", pid, err)
	}
	return parseEnviron(data), nil
}

// ReadSelf returns the environment variables of the current process.
func ReadSelf() map[string]string {
	return parseStrings(os.Environ())
}

// parseEnviron parses the NUL-delimited environ(7) byte slice produced by
// /proc/<pid>/environ into a map.
func parseEnviron(data []byte) map[string]string {
	parts := strings.Split(string(data), "\x00")
	return parseStrings(parts)
}

// parseStrings converts a slice of "KEY=VALUE" strings into a map.
func parseStrings(entries []string) map[string]string {
	env := make(map[string]string, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		idx := strings.IndexByte(entry, '=')
		if idx < 0 {
			// Key with no value — store as empty string.
			env[entry] = ""
			continue
		}
		env[entry[:idx]] = entry[idx+1:]
	}
	return env
}

// ParsePID converts a string to a PID, returning a descriptive error on
// failure.
func ParsePID(s string) (int, error) {
	pid, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || pid <= 0 {
		return 0, fmt.Errorf("invalid PID %q: must be a positive integer", s)
	}
	return pid, nil
}
