package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"plan-balance-service/internal/dto"
	"plan-balance-service/internal/model"
	"plan-balance-service/internal/repository"
)

var (
	ErrCategoryNotFound = errors.New("CATEGORY_NOT_FOUND")
	ErrCategoryExists   = errors.New("CATEGORY_ALREADY_EXISTS")
)

type CategoryService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]dto.CategoryResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	SeedDefaultCategories(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Check if exists
	existing, _ := s.repo.GetByNameAndType(ctx, userID, req.Name, req.Type)
	if existing != nil {
		return nil, ErrCategoryExists
	}

	category := &model.Category{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Type:        req.Type,
		Icon:        req.Icon,
		Color:       req.Color,
		BudgetLimit: req.BudgetLimit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return s.mapToResponse(category), nil
}

func (s *categoryService) GetAll(ctx context.Context, userID uuid.UUID) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		res[i] = *s.mapToResponse(&c)
	}
	return res, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return s.mapToResponse(category), nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Update fields
	category.Name = req.Name
	category.Icon = req.Icon
	category.Color = req.Color
	category.BudgetLimit = req.BudgetLimit
	category.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return s.mapToResponse(category), nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *categoryService) SeedDefaultCategories(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	incomes := []string{"Allowances", "Salary", "Petty Cash", "Bonus", "Others"}
	expenses := []string{"Food", "Social Life", "Pets", "Transportation", "Household", "Beauty", "Health", "Education", "Subscription", "Gift"}

	now := time.Now()

	for _, name := range incomes {
		cat := &model.Category{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      name,
			Type:      model.CategoryIncome,
			Icon:      "trending_up",
			Color:     "#10B981",
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.repo.CreateTx(ctx, tx, cat); err != nil {
			return err
		}
	}

	for _, name := range expenses {
		cat := &model.Category{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      name,
			Type:      model.CategoryExpense,
			Icon:      "trending_down",
			Color:     "#EF4444",
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.repo.CreateTx(ctx, tx, cat); err != nil {
			return err
		}
	}

	return nil
}

func (s *categoryService) mapToResponse(c *model.Category) *dto.CategoryResponse {
	return &dto.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Type:        c.Type,
		Icon:        c.Icon,
		Color:       c.Color,
		BudgetLimit: c.BudgetLimit,
	}
}
