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
	CurrentNotesTool ToolType `json:"currentNotesTool"`
	NotesVaultPath   string   `json:"notesVaultPath"`
	FeedContentDir   string   `json:"feedContentDir"`
}

func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	rdcConfigDir := filepath.Join(base, "rdc")
	if err := os.MkdirAll(rdcConfigDir, 0o755); err != nil {
		return "", err
	}
	return rdcConfigDir, nil
}

func ConfigFilePath() (string, error) {
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
	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}

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

func (c *Config) GetFeedDir() (string, error) {
	if c.FeedContentDir != "" {
		return c.FeedContentDir, nil
	}

	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	// Default path based on known repo structure
	return filepath.Join(dir, "data", "ryan-landing", "src", "src", "content", "feed"), nil
}

func Save(cfg *Config) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
