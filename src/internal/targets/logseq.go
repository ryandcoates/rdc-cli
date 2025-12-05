package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"rdc/internal/config"
	"rdc/internal/domain"
)

type LogseqTarget struct{}

func NewLogseqTarget() *LogseqTarget {
	return NewLogseqTarget()
}

func (l *LogseqTarget) AppendToDaily(cfg *config.Config, t domain.Task, day time.Time) (string, error) {
	if cfg.NotesVaultPath == "" {
		return "", fmt.Errorf("root path is empty; run `rdc config --root <path>` first")
	}

	filename := day.Format("2006_01_02") + ".md"
	path := filepath.Join(cfg.NotesVaultPath, "journals", filename)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	line := fmt.Sprintf("- TODO %s", t.Title)
	if t.Due != nil {
		line += " DEADLINE: <" + t.Due.Format("2006-01-02") + ">"
	}
	if t.Topic != "" {
		line += " #" + t.Topic
	}
	line += "\n"

	if _, err := f.WriteString(line); err != nil {
		return "", err
	}

	return path, nil
}
