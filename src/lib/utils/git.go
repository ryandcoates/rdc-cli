package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func GitAdd(repoPath string, filePath string) error {
	cmd := exec.Command("git", "-C", repoPath, "add", filePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %s (%w)", string(output), err)
	}
	return nil
}

func GitCommit(repoPath string, message string) error {
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	cmd := exec.Command("git", "-C", repoPath, "commit", "-m", message)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %s (%w)", string(output), err)
	}
	return nil
}

func GitPush(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "push")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %s (%w)", string(output), err)
	}
	return nil
}

func GitStatusClean(repoPath string) (bool, error) {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git status failed: %s (%w)", string(output), err)
	}
	return len(output) == 0, nil
}

func GitCheckoutNewBranch(repoPath string, branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	cmd := exec.Command("git", "-C", repoPath, "checkout", "-b", branchName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git checkout -b failed: %s (%w)", string(output), err)
	}
	return nil
}

func GitPushUpstream(repoPath string, branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	cmd := exec.Command("git", "-C", repoPath, "push", "--set-upstream", "origin", branchName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push --set-upstream failed: %s (%w)", string(output), err)
	}
	return nil
}

func GitCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %s (%w)", string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}
