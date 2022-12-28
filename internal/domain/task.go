package domain

import (
	"time"
)

type Task struct {
	ID          int32
	Description string
	Done        bool
	StartDate   time.Time
	DueDate     time.Time
	Deleted     bool
}
