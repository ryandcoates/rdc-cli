package fsutil

import (
	"os"
	"path/filepath"
)

// EnsureDir creates a directory and any necessary parents.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes data to a file, creating the directory path if needed.
func WriteFile(path string, content string) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}
