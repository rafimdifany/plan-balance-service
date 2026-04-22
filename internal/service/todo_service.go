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
	ErrTodoNotFound = errors.New("TODO_NOT_FOUND")
)

type TodoService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) (*dto.TodosListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.TodoResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error)
	PatchStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.PatchTodoStatusRequest) (*dto.TodoResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type todoService struct {
	repo         repository.TodoRepository
	categoryRepo repository.CategoryRepository
}

func NewTodoService(repo repository.TodoRepository, categoryRepo repository.CategoryRepository) TodoService {
	return &todoService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *todoService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTodoRequest) (*dto.TodoResponse, error) {
	// Validate Category
	_, err := s.categoryRepo.GetByID(ctx, req.CategoryID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	todo := &model.Todo{
		ID:          uuid.New(),
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if todo.Status == "" {
		todo.Status = model.TodoStatusTodo
	}
	if todo.Priority == "" {
		todo.Priority = model.TodoPriorityMedium
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, todo.ID, userID)
}

func (s *todoService) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) (*dto.TodosListResponse, error) {
	todos, totalCount, err := s.repo.List(ctx, userID, filters, limit, offset)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TodoResponse, len(todos))
	for i, t := range todos {
		res[i] = *s.mapToResponse(&t)
	}

	return &dto.TodosListResponse{
		Message:    "Todos fetched successfully",
		TotalCount: totalCount,
		Data:       res,
	}, nil
}

func (s *todoService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	return s.mapToResponse(todo), nil
}

func (s *todoService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.UpdateTodoRequest) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	// Validate Category if changed
	if todo.CategoryID != req.CategoryID {
		_, err = s.categoryRepo.GetByID(ctx, req.CategoryID, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		todo.CategoryID = req.CategoryID
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.Status = req.Status
	todo.Priority = req.Priority
	todo.DueDate = req.DueDate
	todo.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, userID)
}

func (s *todoService) PatchStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, req dto.PatchTodoStatusRequest) (*dto.TodoResponse, error) {
	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	todo.Status = req.Status
	todo.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id, userID)
}

func (s *todoService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTodoNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id, userID)
}

func (s *todoService) mapToResponse(t *model.Todo) *dto.TodoResponse {
	return &dto.TodoResponse{
		ID:            t.ID,
		CategoryID:    t.CategoryID,
		CategoryName:  t.CategoryName,
		CategoryIcon:  t.CategoryIcon,
		CategoryColor: t.CategoryColor,
		Title:         t.Title,
		Description:   t.Description,
		Status:        t.Status,
		Priority:      t.Priority,
		DueDate:       t.DueDate,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}
