package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	CreateTx(ctx context.Context, tx pgx.Tx, category *model.Category) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error)
	GetByNameAndType(ctx context.Context, userID uuid.UUID, name string, catType model.CategoryType) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `INSERT INTO categories (id, user_id, name, type, icon, color, budget_limit, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, category.ID, category.UserID, category.Name, category.Type, category.Icon, category.Color, category.BudgetLimit, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *categoryRepository) CreateTx(ctx context.Context, tx pgx.Tx, category *model.Category) error {
	query := `INSERT INTO categories (id, user_id, name, type, icon, color, budget_limit, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.Exec(ctx, query, category.ID, category.UserID, category.Name, category.Type, category.Icon, category.Color, category.BudgetLimit, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *categoryRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	query := `SELECT id, user_id, name, type, icon, color, budget_limit, created_at, updated_at 
			  FROM categories WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Icon, &c.Color, &c.BudgetLimit, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error) {
	query := `SELECT id, user_id, name, type, icon, color, budget_limit, created_at, updated_at 
			  FROM categories WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	var c model.Category
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Icon, &c.Color, &c.BudgetLimit, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) GetByNameAndType(ctx context.Context, userID uuid.UUID, name string, catType model.CategoryType) (*model.Category, error) {
	query := `SELECT id, user_id, name, type, icon, color, budget_limit, created_at, updated_at 
			  FROM categories WHERE user_id = $1 AND name = $2 AND type = $3 AND deleted_at IS NULL`
	var c model.Category
	err := r.db.QueryRow(ctx, query, userID, name, catType).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Icon, &c.Color, &c.BudgetLimit, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `UPDATE categories SET name = $1, icon = $2, color = $3, budget_limit = $4, updated_at = $5 
			  WHERE id = $6 AND user_id = $7 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, category.Name, category.Icon, category.Color, category.BudgetLimit, time.Now(), category.ID, category.UserID)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE categories SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, time.Now(), id, userID)
	return err
}
