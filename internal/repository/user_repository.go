package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	CreateTx(ctx context.Context, tx pgx.Tx, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, email, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *model.User) error {
	query := `INSERT INTO users (id, email, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1`
	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
