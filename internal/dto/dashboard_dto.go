package dto

import (
	"github.com/shopspring/decimal"
)

type DashboardResponse struct {
	Message string `json:"message"`
	Data    struct {
		TotalBalance      decimal.Decimal       `json:"total_balance"`
		MonthlyIncome     decimal.Decimal       `json:"monthly_income"`
		MonthlyExpense    decimal.Decimal       `json:"monthly_expense"`
		RecentTransactions []TransactionResponse `json:"recent_transactions"`
		ActiveGoals        []GoalResponse        `json:"active_goals"`
		PendingTodosCount  int                   `json:"pending_todos_count"`
	} `json:"data"`
}
