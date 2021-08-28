package changelog

import (
	"fmt"
	"io"
	"sort"
	"time"
)

type Release struct {
	Entries []Entry
	Name    string
	Date    time.Time
}

func (r *Release) Write(w io.Writer) (int, error) {
	n := 0
	b, err := fmt.Fprintf(w, "## [%s] - %s\n", r.Name, r.Date.Format("2006-01-02"))
	if err != nil {
		return n, fmt.Errorf("error writing release: %w", err)
	}
	n += b

	cp := make([]Entry, len(r.Entries))
	copy(cp, r.Entries)
	sort.SliceStable(cp, func(a, b int) bool {
		lhs := cp[a]
		rhs := cp[b]
		if lhs.Type != rhs.Type {
			return lhs.Type < rhs.Type
		}
		return false
	})

	const invalidEntry = EntryType(255)
	current := invalidEntry
	for _, ent := range cp {
		if ent.Type != current {
			if current != invalidEntry {
				// starting a new one after printing an existing sectio
				b, err = fmt.Fprintln(w)
				if err != nil {
					return n, fmt.Errorf("error writing newline: %w", err)
				}
				b += n
			}
			current = ent.Type
			b, err = fmt.Fprintf(w, "### %s\n", ent.Type)
			if err != nil {
				return n, fmt.Errorf("error writing entry header: %w", err)
			}
			n += b
		}
		b, err = fmt.Fprintf(w, "- %s\n", ent.Summary)
		if err != nil {
			return n, fmt.Errorf("error writing entry header: %w", err)
		}
		n += b
	}

	b, err = fmt.Fprintln(w)
	if err != nil {
		return n, fmt.Errorf("error writing entry header: %w", err)
	}
	n += b
	return n, nil
}

func (r *Release) isValid() bool {
	return !r.Date.IsZero() && r.Name != ""
}
