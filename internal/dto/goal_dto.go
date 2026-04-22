package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type CreateGoalRequest struct {
	AssetID      *uuid.UUID       `json:"asset_id"`
	CategoryID   *uuid.UUID       `json:"category_id"`
	Name         string           `json:"name" validate:"required"`
	TargetAmount decimal.Decimal  `json:"target_amount" validate:"required"`
	Type         model.GoalType   `json:"type" validate:"required,oneof=SAVINGS BUDGET"`
	Period       model.GoalPeriod `json:"period" validate:"required,oneof=MONTHLY YEARLY ONCE"`
	StartDate    time.Time        `json:"start_date" validate:"required"`
	EndDate      *time.Time       `json:"end_date"`
}

type UpdateGoalRequest struct {
	AssetID      *uuid.UUID       `json:"asset_id"`
	CategoryID   *uuid.UUID       `json:"category_id"`
	Name         string           `json:"name" validate:"required"`
	TargetAmount decimal.Decimal  `json:"target_amount" validate:"required"`
	Type         model.GoalType   `json:"type" validate:"required,oneof=SAVINGS BUDGET"`
	Period       model.GoalPeriod `json:"period" validate:"required,oneof=MONTHLY YEARLY ONCE"`
	StartDate    time.Time        `json:"start_date" validate:"required"`
	EndDate      *time.Time       `json:"end_date"`
	IsActive     bool             `json:"is_active"`
}

type GoalResponse struct {
	ID            uuid.UUID        `json:"id"`
	AssetID       *uuid.UUID       `json:"asset_id,omitempty"`
	AssetName     string           `json:"asset_name,omitempty"`
	CategoryID    *uuid.UUID       `json:"category_id,omitempty"`
	CategoryName  string           `json:"category_name,omitempty"`
	Name          string           `json:"name"`
	TargetAmount  decimal.Decimal  `json:"target_amount"`
	CurrentAmount decimal.Decimal  `json:"current_amount"`
	Progress      decimal.Decimal  `json:"progress"` // Percentage 0-100
	Type          model.GoalType   `json:"type"`
	Period        model.GoalPeriod `json:"period"`
	StartDate     time.Time        `json:"start_date"`
	EndDate       *time.Time       `json:"end_date,omitempty"`
	IsActive      bool             `json:"is_active"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type GoalsListResponse struct {
	Message string         `json:"message"`
	Data    []GoalResponse `json:"data"`
}

type GoalDetailResponse struct {
	Message string       `json:"message"`
	Data    GoalResponse `json:"data"`
}
