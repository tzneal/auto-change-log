package changelog

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"
)

type ChangeLog struct {
	Releases []Release
	Header   string
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

	if c.Header != "" {
		// ensure the header has a new line at the end
		if !strings.HasSuffix(c.Header, "\n") {
			c.Header += "\n"
		}
		b, err = fmt.Fprintf(w, "%s", c.Header)
		if err != nil {
			return n, fmt.Errorf("error writing header: %w", err)
		}
		n += b

	}

	c.sortReleases()

	for _, r := range c.Releases {
		b, err = r.Write(w)
		if err != nil {
			return n, err
		}
		n += b
	}
	return n, nil
}

func (c *ChangeLog) sortReleases() {
	// ensure our releases are sorted newest to oldest
	sort.Slice(c.Releases, func(a, b int) bool {
		return c.Releases[a].Date.After(c.Releases[b].Date)
	})
}

//## [Unreleased] - 2020-05-01
var releaseParser = regexp.MustCompile("^## \\[(.*?)\\]\\s*-\\s*(\\d{4}-\\d{2}-\\d{2})")

func (c *ChangeLog) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	inHeader := false
	var currentRelease *Release
	var currentType = EntryType(255)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.ToUpper(line), "# CHANGELOG") {
			inHeader = true
			continue
		}

		// start of a Releases
		if strings.HasPrefix(line, "## ") {
			inHeader = false
			c.Releases = append(c.Releases, Release{})
			currentRelease = &c.Releases[len(c.Releases)-1]

			matches := releaseParser.FindStringSubmatch(line)
			if len(matches) == 3 {
				captures := matches[1:]
				currentRelease.Name = captures[0]
				t, err := time.ParseInLocation("2006-01-02", captures[1], time.Now().Location())
				if err != nil {
					return fmt.Errorf("error parsing Releases time for '%s': %w", line, err)
				}

				currentRelease.Date = t.UTC()
			} else {
				return fmt.Errorf("error parsing Releases %s", line)
			}
		} else if strings.HasPrefix(line, "### ") {
			typestr := strings.TrimSpace(line[3:])
			if err := currentType.UnmarshalYAML([]byte(typestr)); err != nil {
				return fmt.Errorf("error parsing entry type %s", typestr)
			}
		} else if strings.HasPrefix(line, "- ") {
			currentRelease.Entries = append(currentRelease.Entries, Entry{
				Type:    currentType,
				Summary: line[2:],
			})
		} else if strings.TrimSpace(line) == "" {
			if inHeader {
				c.Header += "\n"
			}
			continue
		} else {
			if inHeader {
				c.Header += line
				c.Header += "\n"
				continue
			}
			return fmt.Errorf("unsupported line", line)
		}
	}

	c.sortReleases()
	return nil
}
