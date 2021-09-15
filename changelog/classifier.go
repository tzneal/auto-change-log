package changelog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type Classifier struct {
	DefaultType     EntryType
	IgnoreMerge     bool
	ClassifyRules   []ClassifyRule
	IssueExtractors []IssueExtractor
	Entries         []Entry
	PathRegexp *regexp.Regexp
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

	if c.PathRegexp != nil {
		patch, err := getPatchWithPreviousCommit(commit)
		include := false
		if err == nil {
			for _, ent := range patch.FilePatches() {
				from, to := ent.Files()
				if from != nil && c.PathRegexp.MatchString(from.Path()) {
					include = true
					break
				}
				if to != nil && c.PathRegexp.MatchString(to.Path()) {
					include = true
					break
				}
			}
		}
		// commit didn't match, so skip it
		if !include {
			return nil
		}
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
	for _, r := range c.ClassifyRules {
		if r.Priority > bestRulePriority && r.Match(commit.Message) {
			bestRulePriority = r.Priority
			e.Type = r.EntryType
		}
	}

	// append references to any issues found in the commit message
	for _, ie := range c.IssueExtractors {
		result := ie.Match(commit.Message)
		if result != "" {
			if ie.Remove {
				e.Summary = ie.FilterSummary(e.Summary)
			}
			e.Summary += " "
			e.Summary += result
		}
	}
	c.Entries = append(c.Entries, e)
	return nil
}

func getPatchWithPreviousCommit(commit *object.Commit) (*object.Patch, error) {
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("error getting commit tree: %w", err)
	}

	parent, err := commit.Parent(0)
	if err != nil {
		return nil, fmt.Errorf("error getting parent: %w", err)
	}
	parentTree, err := parent.Tree()
	if err != nil {
		return nil, fmt.Errorf("error getting parent commit tree: %w", err)
	}
	diff, err :=  parentTree.Diff(commitTree)
	if err != nil {
		return nil, fmt.Errorf("error generating diff: %w", err)
	}
	return diff.Patch()
}
