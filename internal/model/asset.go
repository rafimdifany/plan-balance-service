package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssetType string

const (
	AssetBank    AssetType = "BANK"
	AssetCash    AssetType = "CASH"
	AssetEWallet AssetType = "EWALLET"
)

type Asset struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Name      string          `json:"name"`
	Type      AssetType       `json:"type"`
	Icon      string          `json:"icon"`
	Color     string          `json:"color"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
}
