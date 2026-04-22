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

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Todo, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]model.Todo, int, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type todoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	query := `INSERT INTO todos (id, user_id, category_id, title, description, status, priority, due_date, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, query, todo.ID, todo.UserID, todo.CategoryID, todo.Title, todo.Description, todo.Status, todo.Priority, todo.DueDate, todo.CreatedAt, todo.UpdatedAt)
	return err
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Todo, error) {
	query := `
		SELECT t.id, t.user_id, t.category_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
			   c.name as category_name, c.icon as category_icon, c.color as category_color
		FROM todos t
		JOIN categories c ON t.category_id = c.id
		WHERE t.id = $1 AND t.user_id = $2 AND t.deleted_at IS NULL`
	
	var t model.Todo
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.CategoryID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		&t.CategoryName, &t.CategoryIcon, &t.CategoryColor,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *todoRepository) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]model.Todo, int, error) {
	var whereClauses []string
	var args []interface{}
	argCount := 1

	whereClauses = append(whereClauses, fmt.Sprintf("t.user_id = $%d", argCount))
	args = append(args, userID)
	argCount++

	whereClauses = append(whereClauses, "t.deleted_at IS NULL")

	if val, ok := filters["category_id"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.category_id = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["status"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.status = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["priority"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.priority = $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["start_due_date"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.due_date >= $%d", argCount))
		args = append(args, val)
		argCount++
	}
	if val, ok := filters["end_due_date"]; ok {
		whereClauses = append(whereClauses, fmt.Sprintf("t.due_date <= $%d", argCount))
		args = append(args, val)
		argCount++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM todos t WHERE %s", whereSQL)
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.category_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
			   c.name as category_name, c.icon as category_icon, c.color as category_color
		FROM todos t
		JOIN categories c ON t.category_id = c.id
		WHERE %s
		ORDER BY t.due_date ASC, t.created_at DESC
		LIMIT $%d OFFSET $%d`, whereSQL, argCount, argCount+1)
	
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		err := rows.Scan(
			&t.ID, &t.UserID, &t.CategoryID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
			&t.CategoryName, &t.CategoryIcon, &t.CategoryColor,
		)
		if err != nil {
			return nil, 0, err
		}
		todos = append(todos, t)
	}
	
	return todos, totalCount, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	query := `UPDATE todos SET category_id = $1, title = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = $7 
			  WHERE id = $8 AND user_id = $9 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, todo.CategoryID, todo.Title, todo.Description, todo.Status, todo.Priority, todo.DueDate, time.Now(), todo.ID, todo.UserID)
	return err
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE todos SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, time.Now(), id, userID)
	return err
}
