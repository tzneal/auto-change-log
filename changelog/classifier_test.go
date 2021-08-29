package changelog

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestClassifier_ProcessCommit(t *testing.T) {
	type fields struct {
		DefaultType     EntryType
		IgnoreMerge     bool
		ClassifyRules   []ClassifyRule
		IssueExtractors []IssueExtractor
		Entries         []Entry
	}
	tests := []struct {
		name     string
		fields   fields
		commit   *object.Commit
		expected string
	}{
		{
			name: "extract issue number",
			fields: fields{
				DefaultType:   0,
				IgnoreMerge:   false,
				ClassifyRules: nil,
				IssueExtractors: []IssueExtractor{
					{
						RegExp: `(FOO-\d+)`,
					},
				},
				Entries: nil,
			},
			commit: &object.Commit{
				Message: "do some stuff\nThis also fixes FOO-123",
			},
			expected: "do some stuff FOO-123",
		},
		{
			name: "extract issue number and links it",
			fields: fields{
				DefaultType:   0,
				IgnoreMerge:   false,
				ClassifyRules: nil,
				IssueExtractors: []IssueExtractor{
					{
						RegExp:  `(FOO-\d+)`,
						LinkURL: "http://issue.tracker/$1&q=2",
					},
				},
				Entries: nil,
			},
			commit: &object.Commit{
				Message: "do some stuff\nThis also fixes FOO-123",
			},
			expected: "do some stuff [FOO-123](http://issue.tracker/FOO-123&q=2)",
		},
		{
			name: "extract multiple issue numbers",
			fields: fields{
				DefaultType:   0,
				IgnoreMerge:   false,
				ClassifyRules: nil,
				IssueExtractors: []IssueExtractor{
					{
						RegExp: `(FOO-\d+)`,
					},
				},
				Entries: nil,
			},
			commit: &object.Commit{
				Message: "do some stuff\nThis also fixes FOO-123 and FOO-456",
			},
			expected: "do some stuff FOO-123 FOO-456",
		},
		{
			name: "extract multiple issue number and links them",
			fields: fields{
				DefaultType:   0,
				IgnoreMerge:   false,
				ClassifyRules: nil,
				IssueExtractors: []IssueExtractor{
					{
						RegExp:  `(FOO-\d+)`,
						LinkURL: "http://issue.tracker/$1&q=2",
					},
				},
				Entries: nil,
			},
			commit: &object.Commit{
				Message: "do some stuff\nThis also fixes FOO-123 and FOO-456",
			},
			expected: "do some stuff [FOO-123](http://issue.tracker/FOO-123&q=2) [FOO-456](http://issue.tracker/FOO-456&q=2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Classifier{
				DefaultType:     tt.fields.DefaultType,
				IgnoreMerge:     tt.fields.IgnoreMerge,
				ClassifyRules:   tt.fields.ClassifyRules,
				IssueExtractors: tt.fields.IssueExtractors,
				Entries:         tt.fields.Entries,
			}
			if err := c.ProcessCommit(tt.commit); err != nil {
				t.Fatalf("ProcessCommit() error = %v", err)
			}
			if len(c.Entries) != 1 {
				t.Fatalf("expected a single entry to be created, got %d", len(c.Entries))
			}
			entry := c.Entries[0]
			if entry.Summary != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, entry.Summary)
			}
		})
	}
}
