package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, name string) (*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, id int, name string) (*Category, error)
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, id int) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, name string) (*Category, error) {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, name, created_at`

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

func (r *repository) GetByID(ctx context.Context, id int) (*Category, error) {
	query := `SELECT id, name, created_at FROM categories WHERE id = $1`

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
	query := `SELECT id, name, created_at FROM categories WHERE name = $1`

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

func (r *repository) GetAll(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY name ASC`

	var categories []Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	if categories == nil {
		categories = []Category{}
	}

	return categories, nil
}

func (r *repository) Update(ctx context.Context, id int, name string) (*Category, error) {
	query := `UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name, created_at`

	var category Category
	err := r.db.QueryRowxContext(ctx, query, name, id).StructScan(&category)
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

func (r *repository) Exists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return exists, nil
}
