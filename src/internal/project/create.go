package project

import (
	"fmt"
	"path/filepath"
	"rdc/internal/config"
)

// CreateOptions holds the data needed to create a project
type CreateOptions struct {
	Name      string
	Template  string
	OutputDir string
	DryRun    bool
}

type Scaffolder interface {
	CreateFolders() error
	WriteFiles() error
	RunInit() error
}

func Create(cfg *config.Config, opts CreateOptions) error {
	basePath := "."
	if opts.OutputDir != "" {
		basePath = opts.OutputDir
	} else if cfg.CodeRootPath != "" {
		basePath = cfg.CodeRootPath
	}

	repoRoot := filepath.Join(basePath, opts.Name)

	var s Scaffolder
	switch opts.Template {
	default:
		s = &EmptyScaffolder{RepoRoot: repoRoot, Name: opts.Name}
	}

	if opts.DryRun {
		fmt.Printf("Dry Run Scaffolding:\n")
		fmt.Printf("  Repo Root: %s\n", repoRoot)
		return nil
	}

	// ... actual logic like mkdir, git init, template expansion ...
	if err := s.CreateFolders(); err != nil {
		return err
	}
	if err := s.WriteFiles(); err != nil {
		return err
	}
	return s.RunInit()
}
