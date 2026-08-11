package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	gorm.Model
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      Status     `json:"status" gorm:"default:pending"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	UserID      uint       `json:"user_id" gorm:"not null"`
}

type CreateTaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

type PatchTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *Status    `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}
