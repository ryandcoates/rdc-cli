package project

import (
	"fmt"
	"rdc/internal/config"
	"rdc/internal/project"

	"github.com/spf13/cobra"
)

var (
	templateName string
	outputDir    string
	dryRun       bool
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		opts := project.CreateOptions{
			Name:      projectName,
			Template:  templateName,
			OutputDir: outputDir,
			DryRun:    dryRun,
		}

		if err := project.Create(cfg, opts); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&templateName, "template", "t", "default", "Template to use")
	createCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory (defaults to CodeRoot or ./)")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview actions without applying changes")

	ProjectCmd.AddCommand(createCmd)
}
