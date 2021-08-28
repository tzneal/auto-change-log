package changelog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/tzneal/auto-change-log/cmd/auto-change-log/changelog"
)

func TestChangeLog_EmptyWrite(t *testing.T) {
	cl := changelog.New()
	o := bytes.Buffer{}
	_, err := cl.Write(&o)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	exp := "# Changelog\n"
	got := o.String()
	if got != exp {
		t.Errorf("expected '%s', got '%s'", printable(exp), printable(got))
	}
}

func TestChangeLog_SingleRelease(t *testing.T) {
	cl := changelog.New()
	r := changelog.Release{}
	r.Name = "1.0.0"
	r.Date = time.UnixMilli(1234512345123)

	err := cl.AddRelease(r)
	if err != nil {
		t.Fatalf("expected no errror, got %s", err)
	}

	o := bytes.Buffer{}
	_, err = cl.Write(&o)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	exp := "# Changelog\n## [1.0.0] - 2009-02-13\n"
	got := o.String()
	if got != exp {
		t.Errorf("expected '%s', got '%s'", printable(exp), printable(got))
	}
}

func TestChangeLog_SingleReleaseWithEntries(t *testing.T) {
	cl := changelog.New()
	r := changelog.Release{}
	r.Name = "1.0.0"
	r.Date = time.UnixMilli(1234512345123)

	r.Entries = append(r.Entries,
		changelog.Entry{
			Type:    changelog.FixedEntry,
			Summary: "Bug was resolved",
		})

	r.Entries = append(r.Entries,
		changelog.Entry{
			Type:    changelog.AddedEntry,
			Summary: "New blah blah blah",
		})

	err := cl.AddRelease(r)
	if err != nil {
		t.Fatalf("expected no errror, got %s", err)
	}

	o := bytes.Buffer{}
	_, err = cl.Write(&o)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	exp := `# Changelog
## [1.0.0] - 2009-02-13
### Added
- New blah blah blah

### Fixed
- Bug was resolved

`
	got := o.String()
	if got != exp {
		t.Errorf("expected '%s', got '%s'", printable(exp), printable(got))
	}
}
func printable(got string) string {
	r := strings.NewReplacer(
		"\n", "\\n",
		"\t", "\\t")
	return r.Replace(got)
}
