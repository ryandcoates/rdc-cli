package rdc

import (
	"fmt"
	"os"
	"path/filepath"
	"rdc/internal/config"

	"gopkg.in/yaml.v3"

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
		configFile := filepath.Join(configDir, "repos.yaml")

		// Read existing config if present
		var cfg SyncReposConfig
		if f, err := os.Open(configFile); err == nil {
			defer f.Close()
			yaml.NewDecoder(f).Decode(&cfg)
		} else {
			cfg.SyncRepos = new.ReposToSync()
		}

		// Add new repo if flags provided
		if addRepoURL != "" && addRepoName != "" {
			newRepo := new.RepoConfig{URL: addRepoURL, Name: addRepoName}
			found := false
			for _, r := range cfg.SyncRepos {
				if r.URL == newRepo.URL || r.Name == newRepo.Name {
					found = true
					break
				}
			}
			if !found {
				cfg.SyncRepos = append(cfg.SyncRepos, newRepo)
			}
		}

		f, err := os.Create(configFile)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := yaml.NewEncoder(f)
		if err := enc.Encode(&cfg); err != nil {
			return err
		}
		fmt.Printf("Wrote %d repos to %s\n", len(cfg.SyncRepos), configFile)
		return nil
	},
}

func init() {
	syncrepoCmd.Flags().StringVar(&addRepoURL, "url", "", "URL of repo to add")
	syncrepoCmd.Flags().StringVar(&addRepoName, "name", "", "Name of repo to add")
	// This should be called from root.go
}
