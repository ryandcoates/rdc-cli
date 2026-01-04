package rdc

import (
	"fmt"

	"rdc/cmd/rdc/project"

	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "rdc",
	Short: "rdc is your personal CLI helper",
	Long:  "rdc is a CLI tool to manage personal tasks and integrations (e.g. Obsidian, Logseq).",
}

// Execute is called by main.main().
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("rdc: %w", err)
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&Verbose,
		"verbose",
		"v",
		false,
		"enable verbose logging",
	)

	rootCmd.AddCommand(project.ProjectCmd)
}
