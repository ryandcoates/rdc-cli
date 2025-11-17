package domain

import "time"

type Task struct {
	Title     string
	Due       *time.Time // nil if no due date
	Topic     string
	CreatedAt time.Time
}
