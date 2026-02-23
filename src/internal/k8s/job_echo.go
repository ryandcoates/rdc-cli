package k8s

import (
	"bytes"
	"fmt"
	"text/template"
)

// EchoOptions holds parameters for an ad-hoc echo job.
type EchoOptions struct {
	Name      string
	Namespace string
	Image     string
	Message   string
}

var echoTemplate = template.Must(template.New("echo").Parse(`apiVersion: batch/v1
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

// RenderEchoYAML fills the echo job template and returns YAML bytes.
func RenderEchoYAML(opts EchoOptions) ([]byte, error) {
	data := struct {
		EchoOptions
		QuotedMessage string
	}{
		EchoOptions:   opts,
		QuotedMessage: shellQuote(opts.Message),
	}

	var buf bytes.Buffer
	if err := echoTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render echo yaml: %w", err)
	}
	return buf.Bytes(), nil
}
