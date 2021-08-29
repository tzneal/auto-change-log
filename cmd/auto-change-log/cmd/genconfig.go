package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/tzneal/auto-change-log/cmd/auto-change-log/changelog"
)

var genConfigCmd = &cobra.Command{
	Use:   "gen-config",
	Short: "generates a new config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting CWD: %w", err)
		}
		c := changelog.DefaultConfig()
		f, err := os.OpenFile(path.Join(wd, ".auto-change-log"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}
		if err = c.Write(f); err != nil {
			return err
		}

		return f.Close()
	},
}
