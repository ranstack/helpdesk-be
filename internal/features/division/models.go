package division

import "time"

type Division struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
