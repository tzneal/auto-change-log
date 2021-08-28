package changelog

import (
	"fmt"
	"io"
	"sort"
)

type ChangeLog struct {
	release []Release
}

func New() *ChangeLog {
	return &ChangeLog{}
}

func (c *ChangeLog) Write(w io.Writer) (int, error) {
	n := 0
	b, err := fmt.Fprint(w, "# Changelog\n")
	if err != nil {
		return n, fmt.Errorf("error writing header: %w", err)
	}
	n += b

	for _, r := range c.release {
		b, err = r.Write(w)
		if err != nil {
			return n, err
		}
		n += b
	}
	return n, nil
}

func (c *ChangeLog) AddRelease(r Release) error {
	if !r.isValid() {
		return fmt.Errorf("invalid release")
	}
	c.release = append(c.release, r)
	sort.Slice(c.release, func(a, b int) bool {
		return c.release[a].Date.Before(c.release[b].Date)
	})
	return nil
}
