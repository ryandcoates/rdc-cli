package rdc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"rdc/internal/k8s"
)

var echoImage string

func init() {
	echoCmd := &cobra.Command{
		Use:   "echo <message>",
		Short: "Echo a message inside a pod",
		Long: `Submits a one-shot Kubernetes Job that runs 'echo MESSAGE' inside a container.

Use --dry-run to preview the YAML without applying it.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := jobName
			if name == "" {
				name = k8s.GenerateName("echo")
			}

			opts := k8s.EchoOptions{
				Name:      name,
				Namespace: jobNamespace,
				Image:     echoImage,
				Message:   args[0],
			}

			yaml, err := k8s.RenderEchoYAML(opts)
			if err != nil {
				return err
			}

			if jobDryRun {
				fmt.Fprint(os.Stdout, "# dry-run — YAML that would be applied:\n\n")
				_, err = os.Stdout.Write(yaml)
				return err
			}

			appliedName, err := k8s.Apply(yaml, name, jobNamespace)
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

	echoCmd.Flags().StringVar(&echoImage, "image", "busybox:latest", "Container image to use")

	jobCmd.AddCommand(echoCmd)
}
