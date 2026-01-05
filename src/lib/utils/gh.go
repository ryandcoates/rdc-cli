package utils

import (
	"fmt"
	"os/exec"
)

func GhCreate(repoPath string, branch string) error {
	cmd := exec.Command("gh", "repo", "create", branch, "--confirm", "--public", "--remote=origin")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gh repo create failed: %s (%w)", string(output), err)
	}
	return nil
}

func GhRepoClone(repoPath string, repoURL string) error {
	cmd := exec.Command("gh", "repo", "clone", repoURL, repoPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gh repo clone failed: %s (%w)", string(output), err)
	}

	return nil
}
