package rdc

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"rdc/internal/config"
	"rdc/internal/tasks"
)

var (
	taskDue   string
	taskTopic string
)

func init() {
	taskCmd := &cobra.Command{
		Use:   "task [title]",
		Short: "Create a new task in today's daily log",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]

			var duePtr *time.Time
			if taskDue != "" {
				d, err := time.Parse("20060102", taskDue)
				if err != nil {
					return err
				}
				duePtr = &d
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			path, err := tasks.RegisterTask(cfg, tasks.RegisterTaskOptions{
				Title: title,
				Due:   duePtr,
				Topic: taskTopic,
			})
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Printf("Task written to: %s\n", path)
			}

			return nil
		},
	}

	taskCmd.Flags().StringVar(&taskDue, "due", "", "Due date (YYYYMMDD)")
	taskCmd.Flags().StringVar(&taskTopic, "topic", "", "Topic/tag for the task")

	rootCmd.AddCommand(taskCmd)
}
