package output

import (
	"encoding/json"
	"io"

	"github.com/user/envoy-diff/internal/diff"
)

// jsonEntry is the serialisable representation of a diff.Entry.
type jsonEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

// jsonOutput is the top-level JSON document.
type jsonOutput struct {
	Differences []jsonEntry `json:"differences"`
	Total        int         `json:"total"`
	Changed      int         `json:"changed"`
}

func writeJSON(result []diff.Entry, w io.Writer) (int, error) {
	entries := make([]jsonEntry, 0, len(result))
	changed := 0

	for _, e := range result {
		if e.Status == diff.StatusEqual {
			continue
		}
		je := jsonEntry{
			Key:    e.Key,
			Status: statusString(e.Status),
			ValueA: e.ValueA,
			ValueB: e.ValueB,
		}
		entries = append(entries, je)
		changed++
	}

	out := jsonOutput{
		Differences: entries,
		Total:        len(result),
		Changed:      changed,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return 0, err
	}
	return changed, nil
}

func statusString(s diff.Status) string {
	switch s {
	case diff.StatusOnlyInA:
		return "removed"
	case diff.StatusOnlyInB:
		return "added"
	case diff.StatusChanged:
		return "changed"
	default:
		return "equal"
	}
}
