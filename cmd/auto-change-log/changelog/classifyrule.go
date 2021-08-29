package changelog

import (
	"log"
	"regexp"
)

type ClassifyRule struct {
	RegExp    string
	EntryType EntryType
	Priority  int
	regexp    *regexp.Regexp
}

func (r *ClassifyRule) Compile() error {
	var err error
	r.regexp, err = regexp.Compile(r.RegExp)
	return err
}

func (c *ClassifyRule) Match(s string) bool {
	if c.regexp == nil {
		if err := c.Compile(); err != nil {
			log.Fatalf("error matching: %e", err)
		}
	}
	return c.regexp.MatchString(s)
}
