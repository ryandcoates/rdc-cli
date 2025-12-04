package rdc

import (
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new resources (repos, configs, etc)",
}

func init() {
	newCmd.AddCommand(syncrepoCmd)
}
