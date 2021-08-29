package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/tzneal/auto-change-log/changelog"
)

func init() {
	generateCommand.Flags().StringP("name", "n", "Unreleased", "the release name, this must match the newest release if appending")
	generateCommand.Flags().StringP("since", "s", "", "only look at commits since this time, specified as YYYY-MM-DD or a tag name")
	generateCommand.Flags().StringP("until", "u", "", "only look at commits up to this time, specified as YYYY-MM-DD or a tag name")
	generateCommand.Flags().Bool("overwrite", false, "overwrite the current changelog entries. If not set, append.")
}

var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "writes the change log",
	RunE: func(cmd *cobra.Command, args []string) error {

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting CWD: %w", err)
		}
		cfgFilename := path.Join(wd, ".auto-change-log")
		var cfg *changelog.Config
		if fi, _ := os.Stat(cfgFilename); fi == nil {
			fmt.Println("config not found, using default")
			cfg = changelog.DefaultConfig()
		} else {
			cfg, err = changelog.OpenConfig(cfgFilename)
			if err != nil {
				return err
			}
		}
		repoPath := cfg.RepositoryPath
		if repoPath == "" {
			repoPath = wd
		}
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			return fmt.Errorf("error opening repo: %w", err)
		}
		lo := &git.LogOptions{
			Order: git.LogOrderCommitterTime,
		}
		// parse our flags
		since := cmd.Flags().Lookup("since").Value.String()
		if since != "" {
			sinceTime, err := parseTimeFromTagOrTime(repo, since)
			if err != nil {
				log.Fatalf("error parsing since flag: %s", err)
			}
			lo.Since = &sinceTime
		}

		until := cmd.Flags().Lookup("until").Value.String()
		if until != "" {
			untilTime, err := parseTimeFromTagOrTime(repo, until)
			if err != nil {
				log.Fatalf("error parsing until flag: %s", err)
			}
			lo.Until = &untilTime
		}
		cl := changelog.New()
		r := &changelog.Release{
			Name: cmd.Flags().Lookup("name").Value.String(),
			Date: time.Now(),
		}
		// the date in the release is up until the point we looked at logs
		if lo.Until != nil {
			r.Date = *lo.Until
		}

		newRelease := true
		// are we appending to the changelog?
		if cmd.Flags().Lookup("overwrite").Value.String() != "true" {
			f, err := os.Open(cfg.Filename)
			if err == nil {
				err = cl.Read(f)
				if err != nil {
					log.Fatalf("error reading changelog: %s", err)
				}
			}
			// do we append to the existing release, or start a new one?
			if len(cl.Releases) > 0 && r.Name == cl.Releases[0].Name {
				newRelease = false
				r = &cl.Releases[0]
			}

			// if since is not specified, start at the day of the last release so we don't
			// duplicate log messages.  This still has duplicates if commits occurred on the day
			// of the last release, but we de-dupe entries within a release later
			if since == "" && len(cl.Releases) > 0 {
				nextTime := cl.Releases[0].Date
				lo.Since = &nextTime
			}
		}

		// are we creating a new release?
		if newRelease {
			// these will get sorted later
			cl.Releases = append(cl.Releases, *r)
			r = &cl.Releases[len(cl.Releases)-1]
		}

		logs, err := repo.Log(lo)
		if err != nil {
			log.Fatalf("error retrieving logs: %s", err)
		}

		clf := changelog.NewClassifier()
		clf.DefaultType = cfg.DefaultEntryType
		clf.ClassifyRules = cfg.ClassifyRules
		clf.IgnoreMerge = cfg.IgnoreMerge

		err = logs.ForEach(clf.ProcessCommit)
		if err != nil {
			log.Fatalf("error enumerating logs: %s", err)
		}
		fmt.Printf("writing %d new log entries to release '%s'\n", len(clf.Entries), r.Name)

		// do we need to prepend our entries to the existing release's entries?
		if len(r.Entries) > 0 {
			cp := make([]changelog.Entry, len(r.Entries)+len(clf.Entries))
			copy(cp, clf.Entries)
			copy(cp[len(clf.Entries):], r.Entries)
			r.Entries = cp
			r.Cleanup()
		} else {
			r.Entries = clf.Entries
		}

		f, err := os.OpenFile(cfg.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", cfg.Filename, err)
		}
		if _, err = cl.Write(f); err != nil {
			return fmt.Errorf("error writing changelog: %w", err)
		}
		return f.Close()
	},
}

func parseTimeFromTagOrTime(repo *git.Repository, value string) (time.Time, error) {
	// do we have a tag?
	tag, err := repo.Tag(value)
	if err == nil {
		cmt, err := repo.CommitObject(tag.Hash())
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to find commit for tag %s hash %s: %w", value, tag.Hash(), err)
		}
		return cmt.Committer.When, nil
	} else {
		t, err := time.ParseInLocation("2006-01-02", value, time.Now().Location())
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to find tag %s and unable to parse as a time: %w", value, err)
		}
		return t, nil
	}
}
