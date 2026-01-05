package cmd

import (
	"fmt"
	new2 "rdc/cmd/new"
	"rdc/cmd/validate"
	"rdc/internal/sys"

	"github.com/spf13/cobra"
)

var Verbose bool
var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "rdc",
	Short:   "rdc is your personal CLI helper",
	Long:    "rdc is a CLI tool to manage personal tasks and integrations (e.g. Obsidian, Logseq).",
	Version: version,
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
	rootCmd.SetVersionTemplate(fmt.Sprintf("rdc version {{.Version}} (%s/%s)\n", sys.GetOSType(), sys.GetArch()))

	rootCmd.AddCommand(new2.NewCmd)
	rootCmd.AddCommand(validate.ValidateCmd)
}
