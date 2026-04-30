package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of environment variables as key-value pairs.
type EnvMap map[string]string

// ParseFile reads a .env file and returns an EnvMap.
// It supports KEY=VALUE, KEY="VALUE", and KEY='VALUE' formats.
// Lines starting with '#' and blank lines are ignored.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional "export " prefix
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("%s:%d: invalid line (missing '='): %q", path, lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if key == "" {
			return nil, fmt.Errorf("%s:%d: empty key", path, lineNum)
		}

		val = unquote(val)
		env[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	return env, nil
}

// unquote strips surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
