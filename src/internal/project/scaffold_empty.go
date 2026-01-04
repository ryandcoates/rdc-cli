package project

import (
	"fmt"
	"rdc/lib/fsutil"
)

type EmptyScaffolder struct {
	RepoRoot string
	Name     string
}

func (s *EmptyScaffolder) CreateFolders() error {
	structure := GetDefaultStructure(s.RepoRoot, s.Name)

	for _, folder := range structure.Folders {
		if err := fsutil.EnsureDir(folder); err != nil {
			return err
		}
	}
	return nil
}

func (s *EmptyScaffolder) WriteFiles() error {
	structure := GetDefaultStructure(s.RepoRoot, s.Name)

	for path, content := range structure.Files {
		if err := fsutil.WriteFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func (s *EmptyScaffolder) RunInit() error {
	fmt.Printf("Base structure created for '%s' at %s\n", s.Name, s.RepoRoot)
	return nil
}
