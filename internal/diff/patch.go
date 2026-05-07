package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PatchFormat defines the output format for a patch.
type PatchFormat int

const (
	PatchFormatUnified PatchFormat = iota
	PatchFormatShell
	PatchFormatDotenv
)

// Patch generates a textual patch from a slice of Entry values that can be
// applied to reproduce the "B" side of the diff from the "A" side.
func Patch(entries []Entry, format PatchFormat) string {
	switch format {
	case PatchFormatShell:
		return patchShell(entries)
	case PatchFormatDotenv:
		return patchDotenv(entries)
	default:
		return patchUnified(entries)
	}
}

func patchUnified(entries []Entry) string {
	var sb strings.Builder
	for _, key := range sortedEntryKeys(entries) {
		e := entryByKey(entries, key)
		switch e.Status {
		case StatusAdded:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.ValueB)
		case StatusRemoved:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.ValueA)
		case StatusChanged:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.ValueA)
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.ValueB)
		}
	}
	return sb.String()
}

func patchShell(entries []Entry) string {
	var sb strings.Builder
	for _, key := range sortedEntryKeys(entries) {
		e := entryByKey(entries, key)
		switch e.Status {
		case StatusAdded, StatusChanged:
			fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.ValueB)
		case StatusRemoved:
			fmt.Fprintf(&sb, "unset %s\n", e.Key)
		}
	}
	return sb.String()
}

func patchDotenv(entries []Entry) string {
	var sb strings.Builder
	for _, key := range sortedEntryKeys(entries) {
		e := entryByKey(entries, key)
		switch e.Status {
		case StatusAdded, StatusChanged:
			fmt.Fprintf(&sb, "%s=%q\n", e.Key, e.ValueB)
		case StatusRemoved:
			fmt.Fprintf(&sb, "# removed: %s\n", e.Key)
		}
	}
	return sb.String()
}

func sortedEntryKeys(entries []Entry) []string {
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Status != StatusUnchanged {
			keys = append(keys, e.Key)
		}
	}
	sort.Strings(keys)
	return keys
}

func entryByKey(entries []Entry, key string) Entry {
	for _, e := range entries {
		if e.Key == key {
			return e
		}
	}
	return Entry{}
}
