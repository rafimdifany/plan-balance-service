package dto

import (
	"time"

	"github.com/google/uuid"
	"plan-balance-service/internal/model"
)

type CreateTodoRequest struct {
	CategoryID  uuid.UUID          `json:"category_id" validate:"required"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description"`
	Status      model.TodoStatus   `json:"status" validate:"omitempty,oneof=TODO IN_PROGRESS DONE"`
	Priority    model.TodoPriority `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH"`
	DueDate     *time.Time         `json:"due_date"`
}

type UpdateTodoRequest struct {
	CategoryID  uuid.UUID          `json:"category_id" validate:"required"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description"`
	Status      model.TodoStatus   `json:"status" validate:"required,oneof=TODO IN_PROGRESS DONE"`
	Priority    model.TodoPriority `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH"`
	DueDate     *time.Time         `json:"due_date"`
}

type PatchTodoStatusRequest struct {
	Status model.TodoStatus `json:"status" validate:"required,oneof=TODO IN_PROGRESS DONE"`
}

type TodoResponse struct {
	ID            uuid.UUID          `json:"id"`
	CategoryID    uuid.UUID          `json:"category_id"`
	CategoryName  string             `json:"category_name"`
	CategoryIcon  string             `json:"category_icon"`
	CategoryColor string             `json:"category_color"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Status        model.TodoStatus   `json:"status"`
	Priority      model.TodoPriority `json:"priority"`
	DueDate       *time.Time         `json:"due_date"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type TodosListResponse struct {
	Message    string         `json:"message"`
	TotalCount int            `json:"total_count"`
	Data       []TodoResponse `json:"data"`
}

type TodoDetailResponse struct {
	Message string       `json:"message"`
	Data    TodoResponse `json:"data"`
}
