package k8s

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GenerateName returns a unique job name with an rdc- prefix.
func GenerateName(jobType string) string {
	return fmt.Sprintf("rdc-%s-%d", jobType, time.Now().UnixMilli())
}

// Apply pipes pre-rendered YAML to kubectl apply and returns the job name.
func Apply(yaml []byte, name, namespace string) (string, error) {
	if err := requireKubectl(); err != nil {
		return "", err
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("kubectl apply: %w", err)
	}
	return name, nil
}

// WaitAndLogs blocks until the job completes (or fails) then streams logs.
func WaitAndLogs(name, namespace string) error {
	if err := requireKubectl(); err != nil {
		return err
	}

	fmt.Printf("Waiting for job/%s to complete...\n", name)
	wait := exec.Command(
		"kubectl", "wait",
		fmt.Sprintf("job/%s", name),
		"--for=condition=complete",
		"--timeout=120s",
		"-n", namespace,
	)
	wait.Stdout = os.Stdout
	wait.Stderr = os.Stderr
	if err := wait.Run(); err != nil {
		return fmt.Errorf("job did not complete successfully: %w", err)
	}

	fmt.Println("\n--- Pod logs ---")
	logs := exec.Command(
		"kubectl", "logs",
		fmt.Sprintf("job/%s", name),
		"-n", namespace,
	)
	logs.Stdout = os.Stdout
	logs.Stderr = os.Stderr
	return logs.Run()
}

func requireKubectl() error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl not found in PATH — is it installed?")
	}
	return nil
}

// shellQuote wraps s in single quotes, escaping any embedded single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
