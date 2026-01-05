package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rdc/internal/config"
	//"github.com/charmbracelet/bubbletea"

	"gopkg.in/yaml.v3"
)

type RepoConfig struct {
	URL  string
	Name string
}

var reposToSync = []RepoConfig{
	{URL: "https://github.com/ryandcoates/ryan-landing", Name: "ryan-landing"},
	{URL: "https://github.com/ryandcoates/fleet-infra", Name: "fleet-infra"},
}

var configPath string
var configDir string

func init() {
	var err error

	configPath, err = config.ConfigFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving config file path: %v\n", err)
	}

	configDir, err = config.ConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving config directory: %v\n", err)
	}
}

// ReposToSync returns the list of repos to sync.
func ReposToSync() []RepoConfig {
	return reposToSync
}

// EnsureReposSync syncs the given list of repos.
func EnsureReposSync(repos []RepoConfig) error {

	for _, repo := range repos {
		repoPath := filepath.Join(configDir, "data", repo.Name)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
				return fmt.Errorf("failed to create config dir: %w", err)
			}
			cmd := exec.Command("gh", "repo", "clone", repo.URL, repoPath)
			cmd.Stdout = nil
			cmd.Stderr = nil
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to clone repo %s: %w", repo.Name, err)
			}
		} else {
			cmd := exec.Command("git", "-C", repoPath, "pull")
			cmd.Stdout = nil
			cmd.Stderr = nil
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to pull repo %s: %w", repo.Name, err)
			}
		}
	}
	return nil
}

func AddRepoToConfig(configDir string, newRepo RepoConfig) (int, error) {
	configFile := filepath.Join(configDir, "repos.yaml")

	var repos []RepoConfig
	if data, err := os.ReadFile(configFile); err == nil {
		var cfg struct {
			SyncRepos []RepoConfig `yaml:"sync_repos"`
		}
		if err := yaml.Unmarshal(data, &cfg); err == nil {
			repos = cfg.SyncRepos
		}
	}

	// Logic for checking duplicates
	found := false
	for _, r := range repos {
		if r.URL == newRepo.URL || r.Name == newRepo.Name {
			found = true
			break
		}
	}
	if !found {
		repos = append(repos, newRepo)
	}

	// Write back
	f, err := os.Create(configFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	output := struct {
		SyncRepos []RepoConfig `yaml:"sync_repos"`
	}{SyncRepos: repos}

	if err := yaml.NewEncoder(f).Encode(&output); err != nil {
		return 0, err
	}

	return len(repos), nil
}

func main() {

}
