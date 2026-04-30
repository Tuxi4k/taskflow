package task

import "time"

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

type Status string

func (s *Status) IsValid() bool {
	switch *s {
	case StatusTodo, StatusDoing, StatusDone:
		return true
	}
	return false
}

type Task struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	Status    Status    `gorm:"type:varchar(10);not null;default:todo" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
