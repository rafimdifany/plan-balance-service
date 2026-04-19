package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/model"
	"plan-balance-service/internal/repository"
)

var (
	ErrAssetNotFound = errors.New("ASSET_NOT_FOUND")
	ErrAssetExists   = errors.New("ASSET_ALREADY_EXISTS")
)

type AssetService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateAssetRequest) (*dto.AssetResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) (*dto.AssetsListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.AssetResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateAssetRequest) (*dto.AssetResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type assetService struct {
	repo repository.AssetRepository
}

func NewAssetService(repo repository.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateAssetRequest) (*dto.AssetResponse, error) {
	// Check if exists
	existing, _ := s.repo.GetByName(ctx, userID, req.Name)
	if existing != nil {
		return nil, ErrAssetExists
	}

	asset := &model.Asset{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		Color:     req.Color,
		Balance:   req.Balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	return s.mapToResponse(asset), nil
}

func (s *assetService) GetAll(ctx context.Context, userID uuid.UUID) (*dto.AssetsListResponse, error) {
	assets, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalBalance := decimal.NewFromInt(0)
	resData := make([]dto.AssetResponse, len(assets))
	for i, a := range assets {
		resData[i] = *s.mapToResponse(&a)
		totalBalance = totalBalance.Add(a.Balance)
	}

	return &dto.AssetsListResponse{
		Message:      "Assets retrieved successfully",
		TotalBalance: totalBalance,
		Data:         resData,
	}, nil
}

func (s *assetService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.AssetResponse, error) {
	asset, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	return s.mapToResponse(asset), nil
}

func (s *assetService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateAssetRequest) (*dto.AssetResponse, error) {
	asset, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}

	// Update fields
	asset.Name = req.Name
	asset.Icon = req.Icon
	asset.Color = req.Color
	asset.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, asset); err != nil {
		return nil, err
	}

	return s.mapToResponse(asset), nil
}

func (s *assetService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Check if exists
	_, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAssetNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id, userID)
}

func (s *assetService) mapToResponse(a *model.Asset) *dto.AssetResponse {
	return &dto.AssetResponse{
		ID:      a.ID,
		Name:    a.Name,
		Type:    a.Type,
		Icon:    a.Icon,
		Color:   a.Color,
		Balance: a.Balance,
	}
}
