package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/model"
	"plan-balance-service/internal/repository"
)

var (
	ErrTransactionNotFound = errors.New("TRANSACTION_NOT_FOUND")
)

type TransactionService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) (*dto.TransactionsListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.TransactionResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetMonthlySummary(ctx context.Context, userID uuid.UUID, month, year int) (*dto.TransactionSummaryResponse, error)
}

type transactionService struct {
	repo         repository.TransactionRepository
	assetRepo    repository.AssetRepository
	categoryRepo repository.CategoryRepository
	db           *pgxpool.Pool
}

func NewTransactionService(
	repo repository.TransactionRepository,
	assetRepo repository.AssetRepository,
	categoryRepo repository.CategoryRepository,
	db *pgxpool.Pool,
) TransactionService {
	return &transactionService{
		repo:         repo,
		assetRepo:    assetRepo,
		categoryRepo: categoryRepo,
		db:           db,
	}
}

func (s *transactionService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// 1. Validate Asset
	_, err := s.assetRepo.GetByID(ctx, req.AssetID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	// 2. Validate Category
	_, err = s.categoryRepo.GetByID(ctx, req.CategoryID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// 3. Start Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 4. Create Transaction Record
	trans := &model.Transaction{
		ID:         uuid.New(),
		UserID:     userID,
		AssetID:    req.AssetID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Type:       req.Type,
		Note:       req.Note,
		Date:       req.Date,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.CreateTx(ctx, tx, trans); err != nil {
		return nil, err
	}

	// 5. Update Asset Balance
	balanceChange := req.Amount
	if req.Type == model.TransactionExpense {
		balanceChange = balanceChange.Neg()
	}

	if err := s.assetRepo.UpdateBalanceTx(ctx, tx, req.AssetID, balanceChange); err != nil {
		return nil, err
	}

	// 6. Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// 7. Return with names (refetch or manual join)
	return s.GetByID(ctx, trans.ID, userID)
}

func (s *transactionService) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) (*dto.TransactionsListResponse, error) {
	transactions, totalCount, err := s.repo.List(ctx, userID, filters, limit, offset)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TransactionResponse, len(transactions))
	for i, t := range transactions {
		res[i] = *s.mapToResponse(&t)
	}

	return &dto.TransactionsListResponse{
		Message:    "Transactions fetched successfully",
		TotalCount: totalCount,
		Data:       res,
	}, nil
}

func (s *transactionService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.TransactionResponse, error) {
	trans, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return s.mapToResponse(trans), nil
}

func (s *transactionService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	trans, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	// Validate Category if changed
	if trans.CategoryID != req.CategoryID {
		_, err = s.categoryRepo.GetByID(ctx, req.CategoryID, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		trans.CategoryID = req.CategoryID
	}

	trans.Note = req.Note
	trans.Date = req.Date
	trans.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, trans); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, userID)
}

func (s *transactionService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	trans, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTransactionNotFound
		}
		return err
	}

	// Start Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update Asset Balance (Revert)
	balanceChange := trans.Amount
	if trans.Type == model.TransactionIncome {
		balanceChange = balanceChange.Neg()
	}

	if err := s.assetRepo.UpdateBalanceTx(ctx, tx, trans.AssetID, balanceChange); err != nil {
		return err
	}

	// Delete Transaction
	if err := s.repo.DeleteTx(ctx, tx, id, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *transactionService) GetMonthlySummary(ctx context.Context, userID uuid.UUID, month, year int) (*dto.TransactionSummaryResponse, error) {
	income, expense, err := s.repo.GetMonthlySummary(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}

	return &dto.TransactionSummaryResponse{
		Message:      "Summary fetched successfully",
		TotalIncome:  income,
		TotalExpense: expense,
		NetBalance:   income.Sub(expense),
	}, nil
}

func (s *transactionService) mapToResponse(t *model.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:              t.ID,
		AssetID:         t.AssetID,
		AssetName:       t.AssetName,
		CategoryID:      t.CategoryID,
		CategoryName:    t.CategoryName,
		CategoryIcon:    t.CategoryIcon,
		CategoryColor:   t.CategoryColor,
		Type:            t.Type,
		Amount:          t.Amount,
		Note:            t.Note,
		TransactionDate: t.Date,
	}
}
