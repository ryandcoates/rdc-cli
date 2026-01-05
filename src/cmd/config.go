package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"rdc/internal/config"
)

var (
	cfgCurrentNotesTool string
	cfgNotesVaultPath   string
	cfgFeedContentPath  string
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "View or update rdc configuration",
		Long:  "View or set the current target tool (obsidian/logseq) and root path to its vault/store.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			updated := false

			if cfgCurrentNotesTool != "" {
				switch cfgCurrentNotesTool {
				case string(config.ToolObsidian), string(config.ToolLogseq):
					cfg.CurrentNotesTool = config.ToolType(cfgCurrentNotesTool)
					updated = true
				default:
					return fmt.Errorf("invalid tool %q (expected 'obsidian' or 'logseq')", cfgCurrentNotesTool)
				}
			}

			if cfgNotesVaultPath != "" {
				cfg.NotesVaultPath = cfgNotesVaultPath
				updated = true
			}

			if cfgFeedContentPath != "" {
				cfg.FeedContentDir = cfgFeedContentPath
				updated = true
			}

			if updated {
				if err := config.Save(cfg); err != nil {
					return err
				}
				fmt.Println("Configuration updated.")
			} else {
				fmt.Printf("Current Notes tool: %s\n", cfg.CurrentNotesTool)
				fmt.Printf("Notes Vault Path   : %s\n", cfg.NotesVaultPath)

				feedDir, _ := cfg.GetFeedDir()
				fmt.Printf("Feed content dir  : %s (configured: %q)\n", feedDir, cfg.FeedContentDir)
			}

			return nil
		},
	}

	configCmd.Flags().StringVar(&cfgCurrentNotesTool, "tool", "", "Target Notes tool (obsidian or logseq)")
	configCmd.Flags().StringVar(&cfgNotesVaultPath, "root", "", "Path to Notes Vault")
	configCmd.Flags().StringVar(&cfgFeedContentPath, "feed", "", "Directory for feed content")

	rootCmd.AddCommand(configCmd)
}
