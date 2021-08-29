package changelog

import (
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type Classifier struct {
	DefaultType EntryType
	Rules       []ClassifyRule
	IgnoreMerge bool
	Entries     []Entry
}

func NewClassifier() *Classifier {
	clf := &Classifier{
		DefaultType: ChangedEntry,
	}
	return clf
}
func (c *Classifier) ProcessCommit(commit *object.Commit) error {
	if c.IgnoreMerge && strings.HasPrefix(commit.Message, "Merge ") {
		return nil
	}

	e := Entry{
		Type:    c.DefaultType,
		Summary: commit.Message,
	}
	// truncate to first line
	if idx := strings.Index(e.Summary, "\n"); idx != -1 {
		e.Summary = e.Summary[0:idx]
	}

	bestRulePriority := -1
	for _, r := range c.Rules {
		if r.Priority > bestRulePriority && r.Match(commit.Message) {
			bestRulePriority = r.Priority
			e.Type = r.EntryType
		}
	}
	c.Entries = append(c.Entries, e)
	return nil
}
