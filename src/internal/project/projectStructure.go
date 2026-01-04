package project

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates/*
var templateFS embed.FS

type ProjectStructure struct {
	Folders []string
	Files   map[string]string
}

func GetDefaultStructure(repoRoot string, name string) ProjectStructure {

	gitignore, _ := templateFS.ReadFile("templates/common/gitignore")

	return ProjectStructure{
		Folders: []string{
			repoRoot,
			filepath.Join(repoRoot, "src"),
			filepath.Join(repoRoot, "docs"),
			filepath.Join(repoRoot, "docs", "arch"),
			filepath.Join(repoRoot, "inf"),
			filepath.Join(repoRoot, ".github"),
		},
		Files: map[string]string{
			filepath.Join(repoRoot, "README.md"):  fmt.Sprintf("# %s\n", name),
			filepath.Join(repoRoot, ".gitignore"): string(gitignore),
		},
	}
}
