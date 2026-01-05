package new

import (
	"fmt"
	"rdc/internal/config"

	"rdc/internal/new"

	"github.com/spf13/cobra"
)

var (
	addRepoURL  string
	addRepoName string
)

var syncrepoCmd = &cobra.Command{
	Use:   "syncrepo",
	Short: "Add a repo to ~/.config/tnctl/config.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if addRepoURL == "" || addRepoName == "" {
			// No arguments: print help and exit
			return cmd.Help()
		}

		configDir, err := config.ConfigDir()
		if err != nil {
			return err
		}

		newRepo := new.RepoConfig{URL: addRepoURL, Name: addRepoName}
		count, err := new.AddRepoToConfig(configDir, newRepo)
		if err != nil {
			return fmt.Errorf("failed to update repo config: %w", err)
		}

		fmt.Printf("Wrote %d repos to config folder\n", count)
		return nil
	},
}

func init() {
	syncrepoCmd.Flags().StringVar(&addRepoURL, "url", "", "URL of repo to add")
	syncrepoCmd.Flags().StringVar(&addRepoName, "name", "", "Name of repo to add")

	NewCmd.AddCommand(syncrepoCmd)
}
