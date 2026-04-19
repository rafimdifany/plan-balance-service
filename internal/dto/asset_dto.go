package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type AssetResponse struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Type    model.AssetType `json:"type"`
	Icon    string          `json:"icon"`
	Color   string          `json:"color"`
	Balance decimal.Decimal `json:"balance"`
}

type CreateAssetRequest struct {
	Name    string          `json:"name" binding:"required,max=100"`
	Type    model.AssetType `json:"type" binding:"required,oneof=BANK CASH EWALLET"`
	Icon    string          `json:"icon" binding:"required"`
	Color   string          `json:"color" binding:"required"`
	Balance decimal.Decimal `json:"balance"`
}

type UpdateAssetRequest struct {
	Name  string `json:"name" binding:"required,max=100"`
	Icon  string `json:"icon" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type AssetsListResponse struct {
	Message      string          `json:"message"`
	TotalBalance decimal.Decimal `json:"total_balance"`
	Data         []AssetResponse `json:"data"`
}

type AssetDetailResponse struct {
	Message string        `json:"message"`
	Data    AssetResponse `json:"data"`
}
