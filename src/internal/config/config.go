package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type ToolType string

const (
	ToolObsidian ToolType = "obsidian"
	ToolLogseq   ToolType = "logseq"
)

type Config struct {
	CurrentTool ToolType `json:"currentTool"`
	RootPath    string   `json:"rootPath"`
}

func configFilePath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	rdcDir := filepath.Join(base, "rdc")
	if err := os.MkdirAll(rdcDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(rdcDir, "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return an empty default config
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
