package validate

import (
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate environment dependencies and system requirements",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {

}
