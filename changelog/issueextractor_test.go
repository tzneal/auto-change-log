package changelog

import (
	"testing"
)

func TestIssueExtractor_Match(t *testing.T) {
	type fields struct {
		RegExp  string
		LinkURL string
	}

	tests := []struct {
		name   string
		fields fields
		arg    string
		want   string
	}{
		{
			"empty",
			fields{
				RegExp:  "",
				LinkURL: "",
			},
			"",
			"",
		},
		{
			"match only",
			fields{
				RegExp:  `(#\d+)`,
				LinkURL: "",
			},
			"Fixes #123",
			"#123",
		},
		{
			"match and append URL",
			fields{
				RegExp:  `(FOO-\d+)`,
				LinkURL: "http://bug.tracker/id=$1&q=1",
			},
			"Fixes FOO-123",
			"[FOO-123](http://bug.tracker/id=FOO-123&q=1)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IssueExtractor{
				RegExp:  tt.fields.RegExp,
				LinkURL: tt.fields.LinkURL,
			}
			if got := i.Match(tt.arg); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
