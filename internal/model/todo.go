package model

import (
	"time"

	"github.com/google/uuid"
)

type TodoStatus string

const (
	TodoStatusTodo       TodoStatus = "TODO"
	TodoStatusInProgress TodoStatus = "IN_PROGRESS"
	TodoStatusDone       TodoStatus = "DONE"
)

type TodoPriority string

const (
	TodoPriorityLow    TodoPriority = "LOW"
	TodoPriorityMedium TodoPriority = "MEDIUM"
	TodoPriorityHigh   TodoPriority = "HIGH"
)

type Todo struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	CategoryID  uuid.UUID    `json:"category_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`

	// Joined fields
	CategoryName  string `json:"category_name,omitempty"`
	CategoryIcon  string `json:"category_icon,omitempty"`
	CategoryColor string `json:"category_color,omitempty"`
}
