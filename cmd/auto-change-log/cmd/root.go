package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auto-change-log",
	Short: "auto-change-log is a tool for generating CHANGEOG.md based off of git commit history",
}

func init() {
	rootCmd.AddCommand(generateCommand)
	rootCmd.AddCommand(genConfigCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
