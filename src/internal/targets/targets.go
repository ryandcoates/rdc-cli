package targets

import (
	"fmt"
	"time"

	"rdc/internal/config"
	"rdc/internal/domain"
)

type Target interface {
	AppendToDaily(cfg *config.Config, t domain.Task, day time.Time) (string, error)
}

func NewTarget(cfg *config.Config) (Target, error) {
	switch cfg.CurrentTool {
	case config.ToolObsidian:
		return NewObsidianTarget(), nil
	case config.ToolLogseq:
		return NewLogseqTarget(), nil
	default:
		return nil, fmt.Errorf("unknown tool %q", cfg.CurrentTool)
	}
}
