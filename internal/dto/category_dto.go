package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type CategoryResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Type        model.CategoryType `json:"type"`
	Icon        string           `json:"icon"`
	Color       string           `json:"color"`
	BudgetLimit *decimal.Decimal `json:"budget_limit"`
}

type CreateCategoryRequest struct {
	Name        string           `json:"name" binding:"required,max=100"`
	Type        model.CategoryType `json:"type" binding:"required,oneof=INCOME EXPENSE"`
	Icon        string           `json:"icon" binding:"required"`
	Color       string           `json:"color" binding:"required"`
	BudgetLimit *decimal.Decimal `json:"budget_limit"`
}

type UpdateCategoryRequest struct {
	Name        string           `json:"name" binding:"required,max=100"`
	Icon        string           `json:"icon" binding:"required"`
	Color       string           `json:"color" binding:"required"`
	BudgetLimit *decimal.Decimal `json:"budget_limit"`
}

type CategoriesListResponse struct {
	Message string             `json:"message"`
	Data    []CategoryResponse `json:"data"`
}

type CategoryDetailResponse struct {
	Message string           `json:"message"`
	Data    CategoryResponse `json:"data"`
}
