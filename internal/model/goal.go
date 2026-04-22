package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GoalType string

const (
	GoalTypeSavings GoalType = "SAVINGS"
	GoalTypeBudget  GoalType = "BUDGET"
)

type GoalPeriod string

const (
	GoalPeriodMonthly GoalPeriod = "MONTHLY"
	GoalPeriodYearly  GoalPeriod = "YEARLY"
	GoalPeriodOnce    GoalPeriod = "ONCE"
)

type Goal struct {
	ID            uuid.UUID       `json:"id"`
	UserID        uuid.UUID       `json:"user_id"`
	AssetID       *uuid.UUID      `json:"asset_id"`
	CategoryID    *uuid.UUID      `json:"category_id"`
	Name          string          `json:"name"`
	TargetAmount  decimal.Decimal `json:"target_amount"`
	CurrentAmount decimal.Decimal `json:"current_amount"`
	Type          GoalType        `json:"type"`
	Period        GoalPeriod      `json:"period"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       *time.Time      `json:"end_date"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     *time.Time      `json:"deleted_at,omitempty"`

	// Joined/Calculated fields
	AssetName    string `json:"asset_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
}
