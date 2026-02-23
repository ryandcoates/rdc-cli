package k8s

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

// JobOptions holds the parameters for an ad-hoc echo job.
type JobOptions struct {
	Name      string
	Namespace string
	Image     string
	Message   string
}

// jobTemplate is the embedded Kubernetes Job spec template.
// Parameters are substituted via Go's text/template.
var jobTemplate = template.Must(template.New("job").Parse(`apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    app.kubernetes.io/managed-by: rdc
    rdc/job-type: echo
spec:
  ttlSecondsAfterFinished: 120
  backoffLimit: 0
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: echo
          image: {{ .Image }}
          command: ["sh", "-c", "echo {{ .QuotedMessage }}"]
`))

// GenerateName returns a unique job name with an rdc- prefix.
func GenerateName() string {
	return fmt.Sprintf("rdc-echo-%d", time.Now().UnixMilli())
}

// RenderJobYAML fills the job template with opts and returns YAML bytes.
func RenderJobYAML(opts JobOptions) ([]byte, error) {
	data := struct {
		JobOptions
		QuotedMessage string
	}{
		JobOptions:    opts,
		QuotedMessage: shellQuote(opts.Message),
	}

	var buf bytes.Buffer
	if err := jobTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render job yaml: %w", err)
	}
	return buf.Bytes(), nil
}

// Apply submits the job to the cluster via kubectl apply and returns the
// applied job name.
func Apply(opts JobOptions) (string, error) {
	if err := requireKubectl(); err != nil {
		return "", err
	}

	yaml, err := RenderJobYAML(opts)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("kubectl apply: %w", err)
	}
	return opts.Name, nil
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

// DryRun writes the rendered YAML to w without submitting it.
func DryRun(opts JobOptions, w io.Writer) error {
	yaml, err := RenderJobYAML(opts)
	if err != nil {
		return err
	}
	_, err = w.Write(yaml)
	return err
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
