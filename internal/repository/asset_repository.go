package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/model"
)

type AssetRepository interface {
	Create(ctx context.Context, asset *model.Asset) error
	CreateTx(ctx context.Context, tx pgx.Tx, asset *model.Asset) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Asset, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Asset, error)
	GetByName(ctx context.Context, userID uuid.UUID, name string) (*model.Asset, error)
	Update(ctx context.Context, asset *model.Asset) error
	UpdateBalanceTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, amount decimal.Decimal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type assetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, asset *model.Asset) error {
	query := `INSERT INTO assets (id, user_id, name, type, icon, color, balance, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, asset.ID, asset.UserID, asset.Name, asset.Type, asset.Icon, asset.Color, asset.Balance, asset.CreatedAt, asset.UpdatedAt)
	return err
}

func (r *assetRepository) CreateTx(ctx context.Context, tx pgx.Tx, asset *model.Asset) error {
	query := `INSERT INTO assets (id, user_id, name, type, icon, color, balance, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.Exec(ctx, query, asset.ID, asset.UserID, asset.Name, asset.Type, asset.Icon, asset.Color, asset.Balance, asset.CreatedAt, asset.UpdatedAt)
	return err
}

func (r *assetRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Asset, error) {
	query := `SELECT id, user_id, name, type, icon, color, balance, created_at, updated_at 
			  FROM assets WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []model.Asset
	for rows.Next() {
		var a model.Asset
		err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Icon, &a.Color, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *assetRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Asset, error) {
	query := `SELECT id, user_id, name, type, icon, color, balance, created_at, updated_at 
			  FROM assets WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	var a model.Asset
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Icon, &a.Color, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *assetRepository) GetByName(ctx context.Context, userID uuid.UUID, name string) (*model.Asset, error) {
	query := `SELECT id, user_id, name, type, icon, color, balance, created_at, updated_at 
			  FROM assets WHERE user_id = $1 AND name = $2 AND deleted_at IS NULL`
	var a model.Asset
	err := r.db.QueryRow(ctx, query, userID, name).Scan(&a.ID, &a.UserID, &a.Name, &a.Type, &a.Icon, &a.Color, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *assetRepository) Update(ctx context.Context, asset *model.Asset) error {
	query := `UPDATE assets SET name = $1, icon = $2, color = $3, updated_at = $4 
			  WHERE id = $5 AND user_id = $6 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, asset.Name, asset.Icon, asset.Color, time.Now(), asset.ID, asset.UserID)
	return err
}

func (r *assetRepository) UpdateBalanceTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, amount decimal.Decimal) error {
	query := `UPDATE assets SET balance = balance + $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := tx.Exec(ctx, query, amount, time.Now(), id)
	return err
}

func (r *assetRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE assets SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, time.Now(), id, userID)
	return err
}
