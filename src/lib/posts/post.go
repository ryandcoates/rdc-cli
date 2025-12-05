package posts

import (
	"fmt"
	"os"
	"path/filepath"
	"rdc/internal/config"
	"rdc/internal/domain"
	"rdc/lib/utils"
	"time"
)

type RegisterPostOptions struct {
	Title    string
	DateTime *time.Time
	Content  string
}

var configDir string

func init() {
	var err error

	configDir, err = config.ConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving config directory: %v\n", err)
	}
}

func RegisterPost(cfg *config.Config, opts RegisterPostOptions) (string, error) {
	p := domain.Post{
		Title:     opts.Title,
		Content:   opts.Content,
		CreatedAt: time.Now(),
	}

	if opts.DateTime != nil {
		p.DateTime = opts.DateTime
	} else {
		p.DateTime = &p.CreatedAt
	}

	feedDir, err := cfg.GetFeedDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(feedDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create feed directory: %w", err)
	}

	branchName := fmt.Sprintf("feed/%s", p.CreatedAt.Format("20060102-150405"))
	if err := utils.GitCheckoutNewBranch(feedDir, branchName); err != nil {
		return "", fmt.Errorf("failed to create git branch: %w", err)
	}

	filename := p.CreatedAt.Format("20060102-150405") + ".md"
	fullPath := filepath.Join(feedDir, filename)

	fileContent := fmt.Sprintf("---\ndate: %s\ntype: feed\n---\n\n%s\n",
		p.CreatedAt.Format(time.RFC3339),
		p.Content,
	)

	if err := os.WriteFile(fullPath, []byte(fileContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write feed file: %w", err)
	}

	if err := utils.GitAdd(feedDir, filename); err != nil {
		return fullPath, fmt.Errorf("failed to git add file: %w", err)
	}

	commitMsg := fmt.Sprintf("feat(feed): add post %s", filename)
	if opts.Title != "" {
		commitMsg = fmt.Sprintf("feat(feed): %s", opts.Title)
	}
	if err := utils.GitCommit(feedDir, commitMsg); err != nil {
		return fullPath, fmt.Errorf("failed to git commit: %w", err)
	}

	if err := utils.GitPushUpstream(feedDir, branchName); err != nil {
		return fullPath, fmt.Errorf("failed to git push upstream: %w", err)
	}

	return fullPath, nil
}
