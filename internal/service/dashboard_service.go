package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/repository"
)

type DashboardService interface {
	GetSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	assetRepo       repository.AssetRepository
	transactionRepo repository.TransactionRepository
	goalRepo        repository.GoalRepository
	todoRepo        repository.TodoRepository
	goalSvc         GoalService // To use calculateProgress logic
	transSvc        TransactionService
}

func NewDashboardService(
	assetRepo repository.AssetRepository,
	transactionRepo repository.TransactionRepository,
	goalRepo repository.GoalRepository,
	todoRepo repository.TodoRepository,
	goalSvc GoalService,
	transSvc TransactionService,
) DashboardService {
	return &dashboardService{
		assetRepo:       assetRepo,
		transactionRepo: transactionRepo,
		goalRepo:        goalRepo,
		todoRepo:        todoRepo,
		goalSvc:         goalSvc,
		transSvc:        transSvc,
	}
}

func (s *dashboardService) GetSummary(ctx context.Context, userID uuid.UUID) (*dto.DashboardResponse, error) {
	// 1. Total Balance
	assets, err := s.assetRepo.GetAllByUserID(ctx, userID)
	totalBalance := decimal.NewFromInt(0)
	if err == nil {
		for _, a := range assets {
			totalBalance = totalBalance.Add(a.Balance)
		}
	}

	// 2. Monthly Summary
	now := time.Now()
	income, expense, _ := s.transactionRepo.GetMonthlySummary(ctx, userID, int(now.Month()), now.Year())

	// 3. Recent Transactions
	transList, _ := s.transSvc.List(ctx, userID, nil, 5, 0)

	// 4. Active Goals
	goals, _ := s.goalRepo.List(ctx, userID, map[string]interface{}{"is_active": true})
	activeGoals := make([]dto.GoalResponse, 0)
	for _, g := range goals {
		activeGoals = append(activeGoals, *s.goalSvc.(*goalService).calculateProgress(ctx, userID, &g))
	}

	// 5. Pending Todos Count
	_, todosCount, _ := s.todoRepo.List(ctx, userID, map[string]interface{}{"status": "TODO"}, 0, 0)

	res := &dto.DashboardResponse{}
	res.Message = "Dashboard summary fetched successfully"
	res.Data.TotalBalance = totalBalance
	res.Data.MonthlyIncome = income
	res.Data.MonthlyExpense = expense
	if transList != nil {
		res.Data.RecentTransactions = transList.Data
	}
	res.Data.ActiveGoals = activeGoals
	res.Data.PendingTodosCount = todosCount

	return res, nil
}
