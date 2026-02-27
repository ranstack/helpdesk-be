package division

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
	GetAll(ctx context.Context, filter *DivisionListFilter) ([]Division, int, error)
	GetByID(ctx context.Context, id int) (*Division, error)
	GetByName(ctx context.Context, name string) (*Division, error)
	Exists(ctx context.Context, id int) (bool, error)
	Create(ctx context.Context, name string) (*Division, error)
	Update(ctx context.Context, id int, name string, isActive bool) (*Division, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, filter *DivisionListFilter) ([]Division, int, error) {
	whereClause, args := buildDivisionFilterWhereClause(filter)

	countQuery := `SELECT COUNT(*) FROM divisions` + whereClause
	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count divisions: %w", err)
	}

	limitPlaceholder := len(args) + 1
	offsetPlaceholder := len(args) + 2
	query := fmt.Sprintf(`SELECT id, name, is_active, created_at FROM divisions%s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, whereClause, limitPlaceholder, offsetPlaceholder)
	listArgs := append(args, filter.Limit, filter.Offset)

	var divisions []Division
	err := r.db.SelectContext(ctx, &divisions, query, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get divisions: %w", err)
	}

	if divisions == nil {
		divisions = []Division{}
	}

	return divisions, totalItems, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*Division, error) {
	query := `SELECT id, name, is_active, created_at FROM divisions WHERE id = $1`

	var division Division
	err := r.db.GetContext(ctx, &division, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get division: %w", err)
	}

	return &division, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*Division, error) {
	query := `SELECT id, name, is_active, created_at FROM divisions WHERE LOWER(name) = LOWER($1)`

	var division Division
	err := r.db.GetContext(ctx, &division, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get division: %w", err)
	}

	return &division, nil
}

func (r *repository) Exists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM divisions WHERE id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check division existence: %w", err)
	}

	return exists, nil
}

func (r *repository) Create(ctx context.Context, name string) (*Division, error) {
	query := `INSERT INTO divisions (name) VALUES ($1) RETURNING id, name, is_active, created_at`

	var division Division
	err := r.db.QueryRowxContext(ctx, query, name).StructScan(&division)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("division with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to create division: %w", err)
	}

	return &division, nil
}

func (r *repository) Update(ctx context.Context, id int, name string, isActive bool) (*Division, error) {
	query := `UPDATE divisions SET name = $1, is_active = $2 WHERE id = $3 RETURNING id, name, is_active, created_at`

	var division Division
	err := r.db.QueryRowxContext(ctx, query, name, isActive, id).StructScan(&division)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("division with name '%s' already exists", name)
		}
		return nil, fmt.Errorf("failed to update division: %w", err)
	}

	return &division, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM divisions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete division: %w", err)
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

func buildDivisionFilterWhereClause(filter *DivisionListFilter) (string, []interface{}) {
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
