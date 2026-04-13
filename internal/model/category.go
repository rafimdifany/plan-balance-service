package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CategoryType string

const (
	CategoryIncome  CategoryType = "INCOME"
	CategoryExpense CategoryType = "EXPENSE"
)

type Category struct {
	ID          uuid.UUID        `json:"id"`
	UserID      uuid.UUID        `json:"user_id"`
	Name        string           `json:"name"`
	Type        CategoryType     `json:"type"`
	Icon        string           `json:"icon"`
	Color       string           `json:"color"`
	BudgetLimit *decimal.Decimal `json:"budget_limit"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}
