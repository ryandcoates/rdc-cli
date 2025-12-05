package domain

import "time"

type Post struct {
	Title     string
	DateTime  *time.Time // nil if no due date
	Content   string
	CreatedAt time.Time
}
