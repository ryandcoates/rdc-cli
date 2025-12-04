package validate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	//"fmt"
	"gopkg.in/yaml.v3"
)

type ToolStatus struct {
	Name       string
	Host       string
	Found      bool
	Validated  bool
	Version    string
	ErrMessage string
}

func CheckTool(name string) ToolStatus {
	path, err := exec.LookPath(name)
	if err != nil {
		return ToolStatus{
			Name:       name,
			Found:      false,
			Validated:  false,
			ErrMessage: "not found in PATH",
		}
	}

	versionOutput, err := exec.Command(path, "--version").CombinedOutput()
	version := "unknown"
	if err == nil {
		version = parseVersion(string(versionOutput), name)
	}

	return ToolStatus{
		Name:      name,
		Found:     true,
		Validated: true,
		Version:   version,
	}
}

func parseVersion(raw string, tool string) string {
	raw = strings.TrimSpace(raw)
	switch tool {
	case "git":
		parts := strings.Fields(raw)
		if len(parts) >= 3 {
			return parts[2]
		}
	case "gh":
		parts := strings.Fields(raw)
		if len(parts) >= 3 {
			return parts[2]
		}
	}

	return raw
}

func CheckAllGitHubAuth() []ToolStatus {
	home, err := os.UserHomeDir()
	if err != nil {
		return []ToolStatus{{
			Name:       "gh-auth",
			Found:      false,
			Validated:  false,
			ErrMessage: "unable to determine home directory",
		}}
	}

	// Check potential locations for the config file
	paths := []string{
		filepath.Join(home, ".config", "gh", "hosts.yml"),
	}
	if appData := os.Getenv("APPDATA"); appData != "" {
		paths = append(paths, filepath.Join(appData, "GitHub CLI", "hosts.yml"))
	}

	var data []byte
	found := false

	// Try to read from each path until successful
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		return []ToolStatus{{
			Name:       "gh-auth",
			Found:      false,
			Validated:  false,
			ErrMessage: "no GitHub CLI auth config found",
		}}
	}

	var raw map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return []ToolStatus{{
			Name:       "gh-auth",
			Found:      false,
			Validated:  false,
			ErrMessage: "failed to parse gh hosts.yml",
		}}
	}

	var results []ToolStatus
	for host, info := range raw {
		username := "unknown"
		if user, ok := info["user"].(string); ok {
			username = user
		}

		cmd := exec.Command("gh", "auth", "status", "--hostname", host)
		err := cmd.Run()

		results = append(results, ToolStatus{
			Name:      "gh-auth",
			Host:      host,
			Found:     true,
			Validated: err == nil,
			Version:   username,
			ErrMessage: func() string {
				if err != nil {
					return "auth check failed"
				}
				return ""
			}(),
		})
	}

	return results
}
