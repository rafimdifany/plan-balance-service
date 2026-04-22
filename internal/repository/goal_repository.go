package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/model"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *model.Goal) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) ([]model.Goal, error)
	Update(ctx context.Context, goal *model.Goal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type goalRepository struct {
	db *pgxpool.Pool
}

func NewGoalRepository(db *pgxpool.Pool) GoalRepository {
	return &goalRepository{db: db}
}

func (r *goalRepository) Create(ctx context.Context, goal *model.Goal) error {
	query := `INSERT INTO goals (id, user_id, asset_id, category_id, name, target_amount, current_amount, type, period, start_date, end_date, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err := r.db.Exec(ctx, query, goal.ID, goal.UserID, goal.AssetID, goal.CategoryID, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Type, goal.Period, goal.StartDate, goal.EndDate, goal.IsActive, goal.CreatedAt, goal.UpdatedAt)
	return err
}

func (r *goalRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Goal, error) {
	query := `
		SELECT g.id, g.user_id, g.asset_id, g.category_id, g.name, g.target_amount, g.current_amount, g.type, g.period, g.start_date, g.end_date, g.is_active, g.created_at, g.updated_at,
			   a.name as asset_name, c.name as category_name
		FROM goals g
		LEFT JOIN assets a ON g.asset_id = a.id
		LEFT JOIN categories c ON g.category_id = c.id
		WHERE g.id = $1 AND g.user_id = $2 AND g.deleted_at IS NULL`
	
	var g model.Goal
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&g.ID, &g.UserID, &g.AssetID, &g.CategoryID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.Type, &g.Period, &g.StartDate, &g.EndDate, &g.IsActive, &g.CreatedAt, &g.UpdatedAt,
		&g.AssetName, &g.CategoryName,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *goalRepository) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) ([]model.Goal, error) {
	var whereClauses []string
	var args []interface{}
	argCount := 1

	whereClauses = append(whereClauses, fmt.Sprintf("g.user_id = $%d", argCount))
	args = append(args, userID)
	argCount++

	whereClauses = append(whereClauses, "g.deleted_at IS NULL")

	if val, ok := filters["is_active"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("g.is_active = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["type"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("g.type = $%d", argCount))
		args = append(args, val)
		argCount++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
		SELECT g.id, g.user_id, g.asset_id, g.category_id, g.name, g.target_amount, g.current_amount, g.type, g.period, g.start_date, g.end_date, g.is_active, g.created_at, g.updated_at,
			   a.name as asset_name, c.name as category_name
		FROM goals g
		LEFT JOIN assets a ON g.asset_id = a.id
		LEFT JOIN categories c ON g.category_id = c.id
		WHERE %s
		ORDER BY g.created_at DESC`, whereSQL)
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []model.Goal
	for rows.Next() {
		var g model.Goal
		err := rows.Scan(
			&g.ID, &g.UserID, &g.AssetID, &g.CategoryID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.Type, &g.Period, &g.StartDate, &g.EndDate, &g.IsActive, &g.CreatedAt, &g.UpdatedAt,
			&g.AssetName, &g.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	
	return goals, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *model.Goal) error {
	query := `UPDATE goals SET asset_id = $1, category_id = $2, name = $3, target_amount = $4, current_amount = $5, type = $6, period = $7, start_date = $8, end_date = $9, is_active = $10, updated_at = $11 
			  WHERE id = $12 AND user_id = $13 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, goal.AssetID, goal.CategoryID, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Type, goal.Period, goal.StartDate, goal.EndDate, goal.IsActive, time.Now(), goal.ID, goal.UserID)
	return err
}

func (r *goalRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE goals SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, time.Now(), id, userID)
	return err
}
