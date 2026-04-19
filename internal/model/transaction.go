package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionIncome  TransactionType = "INCOME"
	TransactionExpense TransactionType = "EXPENSE"
)

type Transaction struct {
	ID         uuid.UUID        `json:"id"`
	UserID     uuid.UUID        `json:"user_id"`
	AssetID    uuid.UUID        `json:"asset_id"`
	CategoryID uuid.UUID        `json:"category_id"`
	Amount     decimal.Decimal  `json:"amount"`
	Type       TransactionType  `json:"type"`
	Note       *string          `json:"note"`
	Date       time.Time        `json:"date"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  *time.Time       `json:"deleted_at,omitempty"`

	// Joined fields (optional)
	AssetName     string `json:"asset_name,omitempty"`
	CategoryName  string `json:"category_name,omitempty"`
	CategoryIcon  string `json:"category_icon,omitempty"`
	CategoryColor string `json:"category_color,omitempty"`
}
