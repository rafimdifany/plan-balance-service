package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type TransactionRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, trans *model.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]model.Transaction, int, error)
	GetMonthlySummary(ctx context.Context, userID uuid.UUID, month, year int) (decimal.Decimal, decimal.Decimal, error)
	Update(ctx context.Context, trans *model.Transaction) error
	DeleteTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, userID uuid.UUID) error
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTx(ctx context.Context, tx pgx.Tx, trans *model.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, asset_id, category_id, amount, type, note, date, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := tx.Exec(ctx, query, trans.ID, trans.UserID, trans.AssetID, trans.CategoryID, trans.Amount, trans.Type, trans.Note, trans.Date, trans.CreatedAt, trans.UpdatedAt)
	return err
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.asset_id, t.category_id, t.amount, t.type, t.note, t.date, t.created_at, t.updated_at,
			   a.name as asset_name, c.name as category_name, c.icon as category_icon, c.color as category_color
		FROM transactions t
		JOIN assets a ON t.asset_id = a.id
		JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1 AND t.user_id = $2 AND t.deleted_at IS NULL`
	
	var t model.Transaction
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.AssetID, &t.CategoryID, &t.Amount, &t.Type, &t.Note, &t.Date, &t.CreatedAt, &t.UpdatedAt,
		&t.AssetName, &t.CategoryName, &t.CategoryIcon, &t.CategoryColor,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionRepository) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]model.Transaction, int, error) {
	var whereClauses []string
	var args []interface{}
	argCount := 1

	whereClauses = append(whereClauses, fmt.Sprintf("t.user_id = $%d", argCount))
	args = append(args, userID)
	argCount++

	whereClauses = append(whereClauses, "t.deleted_at IS NULL")

	if val, ok := filters["type"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.type = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["asset_id"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.asset_id = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["category_id"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.category_id = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["start_date"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.date >= $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["end_date"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.date <= $%d", argCount))
		args = append(args, val)
		argCount++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions t WHERE %s", whereSQL)
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.asset_id, t.category_id, t.amount, t.type, t.note, t.date, t.created_at, t.updated_at,
			   a.name as asset_name, c.name as category_name, c.icon as category_icon, c.color as category_color
		FROM transactions t
		JOIN assets a ON t.asset_id = a.id
		JOIN categories c ON t.category_id = c.id
		WHERE %s
		ORDER BY t.date DESC, t.created_at DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argCount, argCount+1)
	
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(
			&t.ID, &t.UserID, &t.AssetID, &t.CategoryID, &t.Amount, &t.Type, &t.Note, &t.Date, &t.CreatedAt, &t.UpdatedAt,
			&t.AssetName, &t.CategoryName, &t.CategoryIcon, &t.CategoryColor,
		)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}
	
	return transactions, totalCount, nil
}

func (r *transactionRepository) GetMonthlySummary(ctx context.Context, userID uuid.UUID, month, year int) (decimal.Decimal, decimal.Decimal, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as total_expense
		FROM transactions 
		WHERE user_id = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(YEAR FROM date) = $3 AND deleted_at IS NULL`
	
	var income, expense decimal.Decimal
	err := r.db.QueryRow(ctx, query, userID, month, year).Scan(&income, &expense)
	return income, expense, err
}

func (r *transactionRepository) Update(ctx context.Context, trans *model.Transaction) error {
	query := `UPDATE transactions SET category_id = $1, note = $2, date = $3, updated_at = $4 
			  WHERE id = $5 AND user_id = $6 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, trans.CategoryID, trans.Note, trans.Date, time.Now(), trans.ID, trans.UserID)
	return err
}

func (r *transactionRepository) DeleteTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE transactions SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	_, err := tx.Exec(ctx, query, time.Now(), id, userID)
	return err
}
