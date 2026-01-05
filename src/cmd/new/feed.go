package new

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"rdc/internal/config"
	"rdc/lib/posts"
)

var (
	feedDate  string
	feedTitle string
)

func init() {
	feedCmd := &cobra.Command{
		Use:   "feed [content]",
		Short: "Create a new feed item",
		Long:  "Create a new timestamped markdown file in the feed directory.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content := strings.Join(args, " ")

			var datePtr *time.Time
			if feedDate != "" {
				d, err := time.Parse("20060102", feedDate)
				if err != nil {
					return fmt.Errorf("invalid date formate (expected YYYYMMDD): %w", err)
				}
				datePtr = &d
			}
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			path, err := posts.RegisterPost(cfg, posts.RegisterPostOptions{
				Title:    feedTitle,
				DateTime: datePtr,
				Content:  content,
			})
			if err != nil {
				return err
			}

			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				fmt.Printf("Feed item written to: %s\n", path)
			} else {
				fmt.Println(path)
			}

			return nil
		},
	}

	feedCmd.Flags().StringVar(&feedDate, "date", "", "Date override (YYYYMMDD)")
	feedCmd.Flags().StringVar(&feedTitle, "title", "", "Title for the feed item")

	NewCmd.AddCommand(feedCmd)
}
