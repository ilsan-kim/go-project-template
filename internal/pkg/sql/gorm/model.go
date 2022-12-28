package gorm

import "time"

type TaskDAO struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	StartDate   time.Time
	DueDate     time.Time
	Done        bool
	Deleted     bool
}

func (TaskDAO) TableName() string {
	return "tasks_gorm"
}

type TaskLogDAO struct {
	ID        uint `gorm:"primaryKey"`
	Action    string
	CreatedAt time.Time
}
