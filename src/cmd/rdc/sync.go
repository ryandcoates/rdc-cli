package rdc

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"rdc/internal/new"
)

// spinnerModel is used for the Bubble Tea spinner
type spinnerModel struct {
	spinner spinner.Model
	done    bool
	err     error
}

func (m spinnerModel) Init() tea.Cmd {
	return spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case syncDoneMsg:
		m.err = msg.err
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return "Sync failed: " + m.err.Error() + "\n"
		}
		return "Repository cache is up-to-date.\n"
	}
	return m.spinner.View() + " Syncing repository cache...\n"
}

// SyncReposConfig is used to parse config.yaml
type SyncReposConfig struct {
	SyncRepos []new.RepoConfig `yaml:"syncrepos"`
}

func getReposForSync() ([]new.RepoConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configFile := filepath.Join(home, ".config", "tnctl", "config.yaml")
	var cfg SyncReposConfig
	fileExists := false
	if f, err := os.Open(configFile); err == nil {
		defer f.Close()
		yaml.NewDecoder(f).Decode(&cfg)
		fileExists = true
	}
	if !fileExists || len(cfg.SyncRepos) == 0 {
		// Always create or update config.yaml with built-in repos if missing or empty
		cfg.SyncRepos = new.ReposToSync()
		if err := os.MkdirAll(filepath.Dir(configFile), 0755); err == nil {
			if f, err := os.Create(configFile); err == nil {
				defer f.Close()
				yaml.NewEncoder(f).Encode(&cfg)
			}
		}
	}
	return cfg.SyncRepos, nil
}

func runSyncWithSpinner() error {
	repos, err := getReposForSync()
	if err != nil {
		return err
	}
	m := spinnerModel{spinner: spinner.New()}

	// Run EnsureReposSync in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- new.EnsureReposSync(repos)
	}()

	p := tea.NewProgram(m)
	go func() {
		err := <-done
		p.Send(syncDoneMsg{err: err})
	}()

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// Custom message for sync completion
type syncDoneMsg struct {
	err error
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the backend GitOps repository cache",
	Run: func(cmd *cobra.Command, args []string) {
		err := runSyncWithSpinner()
		if err != nil {
			fmt.Printf("Sync failed: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
