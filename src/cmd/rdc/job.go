package rdc

import (
	"github.com/spf13/cobra"
)

// jobCmd is the parent for all ad-hoc Kubernetes job subcommands.
// Flags shared across every job type live here as persistent flags.
var jobCmd = &cobra.Command{
	Use:   "job <type> [args]",
	Short: "Run ad-hoc jobs on a Kubernetes cluster",
	Long: `Submits one-shot Kubernetes Jobs to the cluster.

Each sub-command corresponds to a job type with its own template and
arguments. Use --dry-run on any sub-command to preview the YAML.`,
}

var (
	jobNamespace string
	jobName      string
	jobWait      bool
	jobDryRun    bool
)

func init() {
	jobCmd.PersistentFlags().StringVarP(&jobNamespace, "namespace", "n", "default", "Kubernetes namespace")
	jobCmd.PersistentFlags().StringVar(&jobName, "name", "", "Job name (default: auto-generated)")
	jobCmd.PersistentFlags().BoolVar(&jobWait, "wait", false, "Wait for the job to complete and stream logs")
	jobCmd.PersistentFlags().BoolVar(&jobDryRun, "dry-run", false, "Print the Job YAML without applying it")

	rootCmd.AddCommand(jobCmd)
}
