package rdc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"rdc/internal/k8s"
)

var podPingImage string

func init() {
	podPingCmd := &cobra.Command{
		Use:   "podping <ip_address>",
		Short: "Ping an IP from inside a pod",
		Long: `Submits a one-shot Kubernetes Job that runs 'ping IP Address' inside a container.

Use --dry-run to preview the YAML without applying it.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := jobName
			if name == "" {
				name = k8s.GenerateName("podping")
			}

			opts := k8s.PodPingOptions{
				Name:      name,
				Namespace: jobNamespace,
				Image:     echoImage,
				IpAddress: args[0],
			}

			yaml, err := k8s.RenderPodPingYAML(opts)
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

	podPingCmd.Flags().StringVar(&podPingImage, "image", "docker.io/jrecord/nettools:latest", "Container image to use")

	jobCmd.AddCommand(podPingCmd)
}
