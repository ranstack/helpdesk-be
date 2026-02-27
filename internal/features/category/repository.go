package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	GetAll(ctx context.Context, filter *CategoryListFilter) ([]Category, int, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
	Exists(ctx context.Context, id int) (bool, error)
	Create(ctx context.Context, name string) (*Category, error)
	Update(ctx context.Context, id int, name string, isActive bool) (*Category, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, filter *CategoryListFilter) ([]Category, int, error) {
	whereClause, args := buildCategoryFilterWhereClause(filter)

	countQuery := `SELECT COUNT(*) FROM categories` + whereClause
	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	limitPlaceholder := len(args) + 1
	offsetPlaceholder := len(args) + 2
	query := fmt.Sprintf(`SELECT id, name, is_active, created_at FROM categories%s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, whereClause, limitPlaceholder, offsetPlaceholder)
	listArgs := append(args, filter.Limit, filter.Offset)

	var categories []Category
	err := r.db.SelectContext(ctx, &categories, query, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	if categories == nil {
		categories = []Category{}
	}

	return categories, totalItems, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*Category, error) {
	query := `SELECT id, name, is_active, created_at FROM categories WHERE id = $1`

	var category Category
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*Category, error) {
	query := `SELECT id, name, is_active, created_at FROM categories WHERE LOWER(name) = LOWER($1)`

	var category Category
	err := r.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *repository) Exists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return exists, nil
}

func (r *repository) Create(ctx context.Context, name string) (*Category, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name, is_active, created_at`

	var category Category
	err := r.db.QueryRowxContext(ctx, query, name).StructScan(&category)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("category with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func (r *repository) Update(ctx context.Context, id int, name string, isActive bool) (*Category, error) {
	query := `UPDATE categories SET name = $1, is_active = $2 WHERE id = $3 RETURNING id, name, is_active, created_at`

	var category Category
	err := r.db.QueryRowxContext(ctx, query, name, isActive, id).StructScan(&category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("category with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func buildCategoryFilterWhereClause(filter *CategoryListFilter) (string, []interface{}) {
	if filter == nil {
		return "", []interface{}{}
	}

	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.Name != "" {
		args = append(args, "%"+filter.Name+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", len(args)))
	}

	if filter.CreatedAt != nil {
		args = append(args, filter.CreatedAt.Format("2006-01-02"))
		conditions = append(conditions, fmt.Sprintf("DATE(created_at) = $%d::date", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}
