package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type TransactionResponse struct {
	ID              uuid.UUID             `json:"id"`
	AssetID         uuid.UUID             `json:"asset_id"`
	AssetName       string                `json:"asset_name"`
	CategoryID      uuid.UUID             `json:"category_id"`
	CategoryName    string                `json:"category_name"`
	CategoryIcon    string                `json:"category_icon"`
	CategoryColor   string                `json:"category_color"`
	Type            model.TransactionType `json:"type"`
	Amount          decimal.Decimal       `json:"amount"`
	Note            *string               `json:"note"`
	TransactionDate time.Time             `json:"transaction_date"`
}

type CreateTransactionRequest struct {
	AssetID    uuid.UUID             `json:"asset_id" binding:"required"`
	CategoryID uuid.UUID             `json:"category_id" binding:"required"`
	Amount     decimal.Decimal       `json:"amount" binding:"required"`
	Type       model.TransactionType `json:"type" binding:"required,oneof=INCOME EXPENSE"`
	Note       *string               `json:"note"`
	Date       time.Time             `json:"date" binding:"required"`
}

type UpdateTransactionRequest struct {
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	Note       *string   `json:"note"`
	Date       time.Time `json:"date" binding:"required"`
}

type TransactionsListResponse struct {
	Message    string                `json:"message"`
	TotalCount int                   `json:"total_count"`
	Data       []TransactionResponse `json:"data"`
}

type TransactionDetailResponse struct {
	Message string              `json:"message"`
	Data    TransactionResponse `json:"data"`
}

type TransactionSummaryResponse struct {
	Message      string          `json:"message"`
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	NetBalance   decimal.Decimal `json:"net_balance"`
}
