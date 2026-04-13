package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"plan-balance-service/internal/config"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/model"
	"plan-balance-service/internal/repository"
	"plan-balance-service/pkg/utils"
)

var (
	ErrEmailExists        = errors.New("EMAIL_ALREADY_EXISTS")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrInvalidToken       = errors.New("INVALID_TOKEN")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponseData, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponseData, error)
	GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (*dto.AuthResponseData, error)
	Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshResponseData, error)
	Logout(ctx context.Context, req dto.RefreshTokenRequest) error
}

type authService struct {
	userRepo    repository.UserRepository
	authRepo    repository.AuthRepository
	sessionRepo repository.SessionRepository
	categorySvc CategoryService
	cfg         *config.Config
	db          *pgxpool.Pool
}

func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
	categorySvc CategoryService,
	cfg *config.Config,
	db *pgxpool.Pool,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
		categorySvc: categorySvc,
		cfg:         cfg,
		db:          db,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponseData, error) {
	// Check if user exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Start Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Create User
	user := &model.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.CreateTx(ctx, tx, user); err != nil {
		return nil, err
	}

	// Create AuthAccount
	auth := &model.AuthAccount{
		ID:           uuid.New(),
		UserID:       user.ID,
		Provider:     model.ProviderEmail,
		PasswordHash: &hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.authRepo.CreateTx(ctx, tx, auth); err != nil {
		return nil, err
	}

	// Seed Default Categories
	if err := s.categorySvc.SeedDefaultCategories(ctx, tx, user.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.generateAuthData(ctx, user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponseData, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	auth, err := s.authRepo.GetByUserIDAndProvider(ctx, user.ID, model.ProviderEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if auth.PasswordHash == nil || !utils.CheckPasswordHash(req.Password, *auth.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.generateAuthData(ctx, user)
}

func (s *authService) GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (*dto.AuthResponseData, error) {
	payload, err := utils.VerifyGoogleToken(ctx, req.IDToken, s.cfg.GoogleClientID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	googleID := payload.Subject

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Create new user if not exists
			user = &model.User{
				ID:        uuid.New(),
				Email:     email,
				Name:      name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Check if this google account is linked
	auth, err := s.authRepo.GetByProviderInfo(ctx, model.ProviderGmail, googleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Link google account
			auth = &model.AuthAccount{
				ID:             uuid.New(),
				UserID:         user.ID,
				Provider:       model.ProviderGmail,
				ProviderUserID: &googleID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := s.authRepo.Create(ctx, auth); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return s.generateAuthData(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshResponseData, error) {
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	session, err := s.sessionRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if session.Revoked || time.Now().After(session.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Revoke old session (Rotation)
	s.sessionRepo.Revoke(ctx, tokenHash)

	// Generate new access token
	accessToken, err := utils.GenerateToken(session.UserID, []byte(s.cfg.JWTSecret), 1*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponseData{
		AccessToken: accessToken,
		ExpiresIn:   3600,
	}, nil
}

func (s *authService) Logout(ctx context.Context, req dto.RefreshTokenRequest) error {
	tokenHash := utils.HashRefreshToken(req.RefreshToken)
	return s.sessionRepo.Revoke(ctx, tokenHash)
}

func (s *authService) generateAuthData(ctx context.Context, user *model.User) (*dto.AuthResponseData, error) {
	accessToken, err := utils.GenerateToken(user.ID, []byte(s.cfg.JWTSecret), 1*time.Hour)
	if err != nil {
		return nil, err
	}

	rawRefreshToken := uuid.New().String()
	tokenHash := utils.HashRefreshToken(rawRefreshToken)

	session := &model.UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		Revoked:          false,
		CreatedAt:        time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &dto.AuthResponseData{
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    3600,
	}, nil
}
