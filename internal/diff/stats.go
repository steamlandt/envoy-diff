package diff

// Stats holds aggregate counts derived from a slice of diff entries.
type Stats struct {
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Total     int
}

// HasDiff reports whether any differences exist.
func (s Stats) HasDiff() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// DiffCount returns the total number of differing entries.
func (s Stats) DiffCount() int {
	return s.Added + s.Removed + s.Changed
}

// Compute derives Stats from a slice of Entry values produced by Compare.
func Compute(entries []Entry) Stats {
	var s Stats
	s.Total = len(entries)
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		default:
			s.Unchanged++
		}
	}
	return s
}

// FormatStats returns a compact human-readable summary line, e.g.:
//   "3 added, 1 removed, 2 changed, 10 unchanged"
// When there are no differences the string "no differences" is returned.
func FormatStats(s Stats) string {
	if !s.HasDiff() {
		return "no differences"
	}
	parts := make([]string, 0, 4)
	if s.Added > 0 {
		parts = append(parts, plural(s.Added, "added"))
	}
	if s.Removed > 0 {
		parts = append(parts, plural(s.Removed, "removed"))
	}
	if s.Changed > 0 {
		parts = append(parts, plural(s.Changed, "changed"))
	}
	if s.Unchanged > 0 {
		parts = append(parts, plural(s.Unchanged, "unchanged"))
	}
	return join(parts)
}

func plural(n int, label string) string {
	if n == 1 {
		return "1 " + label
	}
	return itoa(n) + " " + label
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

func join(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
