package rdc

import (
	"fmt"
	"rdc/internal/validate"

	"github.com/spf13/cobra"
)

var validateToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Validate tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tools called")
		fmt.Println()

		tools := []validate.ToolStatus{
			validate.CheckTool("git"),
			validate.CheckTool("gh"),
		}
		tools = append(tools, validate.CheckAllGitHubAuth()...)

		for _, t := range tools {
			icon := "✅"
			msg := fmt.Sprintf("v: %s", t.Version)

			if t.Name == "gh-auth" {
				msg = fmt.Sprintf("h: %s  u: %s", t.Host, t.Version)
			}

			if !t.Found || !t.Validated {
				icon = "❌"
				if t.ErrMessage != "" {
					msg = t.ErrMessage
				}
			}

			fmt.Printf("%s %s – %s\n", icon, t.Name, msg)
		}
	},
}

func init() {
	validateCmd.AddCommand(validateToolsCmd)
}
