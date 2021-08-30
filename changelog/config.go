package changelog

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Filename         string
	DefaultEntryType EntryType
	RepositoryPath   string
	IgnoreMerge      bool
	ClassifyRules    []ClassifyRule
	IssueExtractor   []IssueExtractor
}

func (c Config) Write(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	return enc.Encode(c)
}

func DefaultConfig() *Config {
	cfg := &Config{
		DefaultEntryType: ChangedEntry,
		IgnoreMerge:      true,
		Filename:         "CHANGELOG.md",
	}
	cfg.IssueExtractor = append(cfg.IssueExtractor, IssueExtractor{
		RegExp:  `(ISSUE-\d+)`,
		LinkURL: "http://my.bug.tracker/$1/view",
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\b(removed|removing)\\b",
		EntryType: RemovedEntry,
		Priority:  50,
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\b(deprecated|deprecating)",
		EntryType: DeprecatedEntry,
		Priority:  50,
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\b(add|added|adding)\\b",
		EntryType: AddedEntry,
		Priority:  50,
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\b(fix|fixed|fixes)\\b",
		EntryType: FixedEntry,
		Priority:  50,
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\bbug\\b",
		EntryType: FixedEntry,
		Priority:  50,
	})
	cfg.ClassifyRules = append(cfg.ClassifyRules, ClassifyRule{
		RegExp:    "(?i)\\bsecurity",
		EntryType: SecurityEntry,
		Priority:  50,
	})
	return cfg
}

func OpenConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(f)
	cfg := &Config{}
	if err = dec.Decode(cfg); err != nil {
		return nil, err
	}

	for _, r := range cfg.ClassifyRules {
		err = r.Compile()
		if err != nil {
			return nil, fmt.Errorf("error compiling classify rule '%s': %w", r.RegExp, err)
		}
	}
	return cfg, nil
}
