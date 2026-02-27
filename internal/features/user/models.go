package user

import "time"

const (
	RoleAdmin = "ADMIN"
	RoleIT    = "IT"
	RoleStaff = "STAFF"
)

var ValidRoles = map[string]bool{
	RoleAdmin: true,
	RoleIT:    true,
	RoleStaff: true,
}

type User struct {
	ID         int       `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Email      string    `db:"email" json:"email"`
	Password   string    `db:"password" json:"-"`
	AvatarURL  *string   `db:"avatar_url" json:"avatarUrl"`
	Phone      *string   `db:"phone" json:"phone"`
	Role       string    `db:"role" json:"role"`
	DivisionID int       `db:"division_id" json:"divisionId"`
	IsActive   bool      `db:"is_active" json:"isActive"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type UserWithDivision struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	Password     string    `db:"password" json:"-"`
	AvatarURL    *string   `db:"avatar_url" json:"avatarUrl"`
	Phone        *string   `db:"phone" json:"phone"`
	Role         string    `db:"role" json:"role"`
	DivisionID   int       `db:"division_id" json:"divisionId"`
	DivisionName string    `db:"division_name" json:"divisionName"`
	IsActive     bool      `db:"is_active" json:"isActive"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
