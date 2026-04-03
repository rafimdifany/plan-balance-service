package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.UserSession) error
	CreateTx(ctx context.Context, tx pgx.Tx, session *model.UserSession) error
	GetByHash(ctx context.Context, hash string) (*model.UserSession, error)
	Revoke(ctx context.Context, hash string) error
	RevokeTx(ctx context.Context, tx pgx.Tx, hash string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *model.UserSession) error {
	query := `INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at, revoked, device_info, ip_addr, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, session.ID, session.UserID, session.RefreshTokenHash, session.ExpiresAt, session.Revoked, session.DeviceInfo, session.IPAddr, session.CreatedAt)
	return err
}

func (r *sessionRepository) CreateTx(ctx context.Context, tx pgx.Tx, session *model.UserSession) error {
	query := `INSERT INTO user_sessions (id, user_id, refresh_token_hash, expires_at, revoked, device_info, ip_addr, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.Exec(ctx, query, session.ID, session.UserID, session.RefreshTokenHash, session.ExpiresAt, session.Revoked, session.DeviceInfo, session.IPAddr, session.CreatedAt)
	return err
}

func (r *sessionRepository) GetByHash(ctx context.Context, hash string) (*model.UserSession, error) {
	query := `SELECT id, user_id, refresh_token_hash, expires_at, revoked, device_info, ip_addr, created_at FROM user_sessions WHERE refresh_token_hash = $1`
	var session model.UserSession
	err := r.db.QueryRow(ctx, query, hash).Scan(&session.ID, &session.UserID, &session.RefreshTokenHash, &session.ExpiresAt, &session.Revoked, &session.DeviceInfo, &session.IPAddr, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, hash string) error {
	query := `UPDATE user_sessions SET revoked = TRUE WHERE refresh_token_hash = $1`
	_, err := r.db.Exec(ctx, query, hash)
	return err
}

func (r *sessionRepository) RevokeTx(ctx context.Context, tx pgx.Tx, hash string) error {
	query := `UPDATE user_sessions SET revoked = TRUE WHERE refresh_token_hash = $1`
	_, err := tx.Exec(ctx, query, hash)
	return err
}

func (r *sessionRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE user_sessions SET revoked = TRUE WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
