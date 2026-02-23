package rdc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"rdc/internal/k8s"
)

var (
	jobNamespace string
	jobName      string
	jobImage     string
	jobWait      bool
	jobDryRun    bool
)

func init() {
	jobCmd := &cobra.Command{
		Use:   "job <message>",
		Short: "Run an ad-hoc echo job on a Kubernetes cluster",
		Long: `Submits a one-shot Kubernetes Job that echoes MESSAGE inside a pod.

The Job spec is generated from a built-in template and parameterised with
the values you supply. Use --dry-run to preview the YAML without applying.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]

			name := jobName
			if name == "" {
				name = k8s.GenerateName()
			}

			opts := k8s.JobOptions{
				Name:      name,
				Namespace: jobNamespace,
				Image:     jobImage,
				Message:   message,
			}

			if jobDryRun {
				fmt.Fprintf(os.Stdout, "# dry-run — YAML that would be applied:\n\n")
				return k8s.DryRun(opts, os.Stdout)
			}

			appliedName, err := k8s.Apply(opts)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Printf("Job %q submitted to namespace %q\n", appliedName, jobNamespace)
			}

			if jobWait {
				return k8s.WaitAndLogs(appliedName, jobNamespace)
			}

			return nil
		},
	}

	jobCmd.Flags().StringVarP(&jobNamespace, "namespace", "n", "default", "Kubernetes namespace")
	jobCmd.Flags().StringVar(&jobName, "name", "", "Job name (default: auto-generated)")
	jobCmd.Flags().StringVar(&jobImage, "image", "busybox:latest", "Container image to use")
	jobCmd.Flags().BoolVar(&jobWait, "wait", false, "Wait for the job to complete and stream logs")
	jobCmd.Flags().BoolVar(&jobDryRun, "dry-run", false, "Print the Job YAML without applying it")

	rootCmd.AddCommand(jobCmd)
}
