package changelog_test

import (
	"testing"

	"github.com/tzneal/auto-change-log/cmd/auto-change-log/changelog"
)

func TestClassifyRule(t *testing.T) {
	rule := changelog.ClassifyRule{
		RegExp:    "(?i)\\badd",
		EntryType: changelog.AddedEntry,
		Priority:  50,
	}

	if rule.Match(" add some stuff") != true {
		t.Errorf("expected match")
	}
	if rule.Match("add some stuff") != true {
		t.Errorf("expected match")
	}

	if rule.Match("badd some stuff") != false {
		t.Errorf("expected no match")
	}
}
