package user

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
	GetAll(ctx context.Context, filter *UserListFilter) ([]UserWithDivision, int, error)
	GetByID(ctx context.Context, id int) (*UserWithDivision, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	Exists(ctx context.Context, id int) (bool, error)
	Create(ctx context.Context, name, email, passwordHash string, avatarURL, phone, role string, divisionID int) (*UserWithDivision, error)
	Update(ctx context.Context, id int, name, phone, role string, divisionID int, isActive bool) (*UserWithDivision, error)
	UpdateAvatar(ctx context.Context, id int, avatarURL string) (*UserWithDivision, error)
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(ctx context.Context, filter *UserListFilter) ([]UserWithDivision, int, error) {
	whereClause, args := buildUserFilterWhereClause(filter)

	countQuery := `SELECT COUNT(*) FROM users` + whereClause
	var totalItems int
	if err := r.db.GetContext(ctx, &totalItems, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	limitPlaceholder := len(args) + 1
	offsetPlaceholder := len(args) + 2
	query := fmt.Sprintf(`
		SELECT u.id, u.name, u.email, u.password, u.avatar_url, u.phone, u.role, u.division_id, d.name as division_name, u.is_active, u.created_at 
		FROM users u 
		INNER JOIN divisions d ON u.division_id = d.id
		%s 
		ORDER BY u.created_at DESC, u.id DESC 
		LIMIT $%d OFFSET $%d
	`, whereClause, limitPlaceholder, offsetPlaceholder)
	listArgs := append(args, filter.Limit, filter.Offset)

	var users []UserWithDivision
	err := r.db.SelectContext(ctx, &users, query, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	if users == nil {
		users = []UserWithDivision{}
	}

	return users, totalItems, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*UserWithDivision, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.avatar_url, u.phone, u.role, u.division_id, d.name as division_name, u.is_active, u.created_at 
		FROM users u 
		INNER JOIN divisions d ON u.division_id = d.id 
		WHERE u.id = $1
	`

	var user UserWithDivision
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, avatar_url, phone, role, division_id, is_active, created_at FROM users WHERE LOWER(email) = LOWER($1)`

	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*User, error) {
	query := `SELECT id, name, email, password, avatar_url, phone, role, division_id, is_active, created_at FROM users WHERE name = $1`

	var user User
	err := r.db.GetContext(ctx, &user, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *repository) Exists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *repository) Create(ctx context.Context, name, email, passwordHash string, avatarURL, phone, role string, divisionID int) (*UserWithDivision, error) {
	query := `
		INSERT INTO users (name, email, password, avatar_url, phone, role, division_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	var userID int
	err := r.db.QueryRowxContext(ctx, query, name, email, passwordHash, avatarURL, phone, role, divisionID).Scan(&userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("user with email '%s' already exists", email)
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return r.GetByID(ctx, userID)
}

func (r *repository) Update(ctx context.Context, id int, name, phone, role string, divisionID int, isActive bool) (*UserWithDivision, error) {
	query := `
		UPDATE users 
		SET name = $1, phone = $2, role = $3, division_id = $4, is_active = $5 
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query, name, phone, role, divisionID, isActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return r.GetByID(ctx, id)
}

func (r *repository) UpdateAvatar(ctx context.Context, id int, avatarURL string) (*UserWithDivision, error) {
	query := `UPDATE users SET avatar_url = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, avatarURL, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update avatar: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	return r.GetByID(ctx, id)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
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

func buildUserFilterWhereClause(filter *UserListFilter) (string, []interface{}) {
	if filter == nil {
		return "", []interface{}{}
	}

	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.Name != "" {
		args = append(args, "%"+filter.Name+"%")
		conditions = append(conditions, fmt.Sprintf("u.name ILIKE $%d", len(args)))
	}

	if filter.Role != "" {
		args = append(args, filter.Role)
		conditions = append(conditions, fmt.Sprintf("u.role = $%d", len(args)))
	}

	if filter.DivisionID > 0 {
		args = append(args, filter.DivisionID)
		conditions = append(conditions, fmt.Sprintf("u.division_id = $%d", len(args)))
	}

	if filter.IsActive != nil {
		args = append(args, *filter.IsActive)
		conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}
