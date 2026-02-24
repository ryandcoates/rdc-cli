package k8s

import (
	"bytes"
	"fmt"
	"text/template"
)

type PodPingOptions struct {
	Name      string
	Namespace string
	Image     string
	IpAddress string
}

var podPingTemplate = template.Must(template.New("podping").Parse(`apiVersion: batch/v1
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
        - name: podping
          image: {{ .Image }}
          command: ["sh", "-c", "ping -c 4 {{ .IpAddress }}"]
`))

func RenderPodPingYAML(opts PodPingOptions) ([]byte, error) {
	data := struct {
		PodPingOptions
		QuotedMessage string
	}{
		PodPingOptions: opts,
		QuotedMessage:  shellQuote(opts.IpAddress),
	}

	var buf bytes.Buffer
	if err := podPingTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render echo yaml: %w", err)
	}
	return buf.Bytes(), nil
}
