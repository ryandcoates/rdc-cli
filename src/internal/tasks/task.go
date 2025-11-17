package tasks

import (
	"time"

	"rdc/internal/config"
	"rdc/internal/domain"
	"rdc/internal/targets"
)

type RegisterTaskOptions struct {
	Title string
	Due   *time.Time
	Topic string
}

func RegisterTask(cfg *config.Config, opts RegisterTaskOptions) (string, error) {
	t := domain.Task{
		Title:     opts.Title,
		Due:       opts.Due,
		Topic:     opts.Topic,
		CreatedAt: time.Now(),
	}

	target, err := targets.NewTarget(cfg)
	if err != nil {
		return "", err
	}

	today := time.Now()
	return target.AppendToDaily(cfg, t, today)
}
