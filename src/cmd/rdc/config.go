package rdc

import (
	"fmt"

	"github.com/spf13/cobra"

	"rdc/internal/config"
)

var (
	cfgTool string
	cfgRoot string
)

func init() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "View or update rdc configuration",
		Long:  "View or set the current target tool (obsidian/logseq) and root path to its vault/store.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load existing config (if any)
			cfg, err := config.Load()
			if err != nil {
				// For first run, you might decide Load() returns a default empty cfg instead of error.
				// Adjust as you like.
				return err
			}

			updated := false

			if cfgTool != "" {
				switch cfgTool {
				case string(config.ToolObsidian), string(config.ToolLogseq):
					cfg.CurrentTool = config.ToolType(cfgTool)
					updated = true
				default:
					return fmt.Errorf("invalid tool %q (expected 'obsidian' or 'logseq')", cfgTool)
				}
			}

			if cfgRoot != "" {
				cfg.RootPath = cfgRoot
				updated = true
			}

			if updated {
				if err := config.Save(cfg); err != nil {
					return err
				}
				fmt.Println("Configuration updated.")
			} else {
				// No flags passed: just show current config
				fmt.Printf("Current tool: %s\n", cfg.CurrentTool)
				fmt.Printf("Root path   : %s\n", cfg.RootPath)
			}

			return nil
		},
	}

	configCmd.Flags().StringVar(&cfgTool, "tool", "", "Target tool (obsidian or logseq)")
	configCmd.Flags().StringVar(&cfgRoot, "root", "", "Root path to vault/graph")

	rootCmd.AddCommand(configCmd)
}
