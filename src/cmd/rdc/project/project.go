package project

import (
	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Work with projects (create, sync, deploy, etc)",
}

func init() {
}
