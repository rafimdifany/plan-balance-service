package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/model"
)

type AuthRepository interface {
	Create(ctx context.Context, auth *model.AuthAccount) error
	CreateTx(ctx context.Context, tx pgx.Tx, auth *model.AuthAccount) error
	GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider model.AuthProvider) (*model.AuthAccount, error)
	GetByProviderInfo(ctx context.Context, provider model.AuthProvider, providerUserID string) (*model.AuthAccount, error)
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(ctx context.Context, auth *model.AuthAccount) error {
	query := `INSERT INTO auth_accounts (id, user_id, provider, provider_user_id, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, auth.ID, auth.UserID, auth.Provider, auth.ProviderUserID, auth.PasswordHash, auth.CreatedAt, auth.UpdatedAt)
	return err
}

func (r *authRepository) CreateTx(ctx context.Context, tx pgx.Tx, auth *model.AuthAccount) error {
	query := `INSERT INTO auth_accounts (id, user_id, provider, provider_user_id, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(ctx, query, auth.ID, auth.UserID, auth.Provider, auth.ProviderUserID, auth.PasswordHash, auth.CreatedAt, auth.UpdatedAt)
	return err
}

func (r *authRepository) GetByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider model.AuthProvider) (*model.AuthAccount, error) {
	query := `SELECT id, user_id, provider, provider_user_id, password_hash, created_at, updated_at FROM auth_accounts WHERE user_id = $1 AND provider = $2`
	var auth model.AuthAccount
	err := r.db.QueryRow(ctx, query, userID, provider).Scan(&auth.ID, &auth.UserID, &auth.Provider, &auth.ProviderUserID, &auth.PasswordHash, &auth.CreatedAt, &auth.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetByProviderInfo(ctx context.Context, provider model.AuthProvider, providerUserID string) (*model.AuthAccount, error) {
	query := `SELECT id, user_id, provider, provider_user_id, password_hash, created_at, updated_at FROM auth_accounts WHERE provider = $1 AND provider_user_id = $2`
	var auth model.AuthAccount
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(&auth.ID, &auth.UserID, &auth.Provider, &auth.ProviderUserID, &auth.PasswordHash, &auth.CreatedAt, &auth.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
