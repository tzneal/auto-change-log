package changelog

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type IssueExtractor struct {
	RegExp  string
	LinkURL string
	Remove  bool
	Label   string
	regexp  *regexp.Regexp
}

func (i *IssueExtractor) Compile() error {
	var err error
	i.regexp, err = regexp.Compile(i.RegExp)
	return err
}

func (i *IssueExtractor) FilterSummary(s string) string {
	if i.regexp == nil {
		if err := i.Compile(); err != nil {
			log.Fatalf("error matching: %e", err)
		}
	}
	return i.regexp.ReplaceAllString(s, "")
}

func (i *IssueExtractor) Match(s string) string {
	if i.regexp == nil {
		if err := i.Compile(); err != nil {
			log.Fatalf("error matching: %e", err)
		}
	}
	matches := i.regexp.FindAllStringSubmatch(s, -1)

	op := strings.Builder{}

	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		if op.Len() > 0 {
			op.WriteByte(' ')
		}
		if i.LinkURL != "" {
			url := strings.Replace(i.LinkURL, "$1", match[1], -1)
			label := match[1]
			if i.Label != "" {
				label = strings.Replace(i.Label, "$1", match[1], -1)
			}
			op.WriteString(fmt.Sprintf("[%s](%s)", label, url))
		} else {
			op.WriteString(match[1])

		}
	}
	return op.String()
}
