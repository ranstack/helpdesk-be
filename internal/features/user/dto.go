package user

import (
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/validator"
	"strings"
	"time"
)

type CreateUserRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	DivisionID int    `json:"divisionId"`
}

type UpdateUserRequest struct {
	Name       string  `json:"name"`
	Phone      *string `json:"phone"`
	Role       string  `json:"role"`
	DivisionID int     `json:"divisionId"`
	IsActive   *bool   `json:"isActive"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatarUrl"`
	Phone     *string   `json:"phone"`
	Role      string    `json:"role"`
	Division  Division  `json:"division"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type Division struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetUsersQuery struct {
	response.PaginationQuery
	Name       string `query:"name"`
	Role       string `query:"role"`
	DivisionID int    `query:"divisionId"`
	IsActive   *bool  `query:"isActive"`
}

type UserListFilter struct {
	Page       int
	Limit      int
	Offset     int
	Name       string
	Role       string
	DivisionID int
	IsActive   *bool
}

func (r *CreateUserRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 50)
	validator.ValidateString(v, "email", r.Email, true, 5, 255)
	if r.Email != "" && !validator.ValidateEmail(r.Email) {
		v.AddError("email", "Must be a valid email address")
	}
	validator.ValidateString(v, "password", r.Password, true, 6, 255)

	role := strings.TrimSpace(r.Role)
	if role == "" {
		v.AddError("role", "Required")
	} else if !ValidRoles[role] {
		v.AddError("role", "Must be one of: ADMIN, IT, STAFF")
	}

	if r.DivisionID <= 0 {
		v.AddError("divisionId", "Required and must be greater than 0")
	}

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (r *UpdateUserRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 50)

	role := strings.TrimSpace(r.Role)
	if role == "" {
		v.AddError("role", "Required")
	} else if !ValidRoles[role] {
		v.AddError("role", "Must be one of: ADMIN, IT, STAFF")
	}

	if r.DivisionID <= 0 {
		v.AddError("divisionId", "Required and must be greater than 0")
	}

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (q *GetUsersQuery) Normalize() (*UserListFilter, error) {
	page, limit, offset := q.NormalizePagination()

	return &UserListFilter{
		Page:       page,
		Limit:      limit,
		Offset:     offset,
		Name:       strings.TrimSpace(q.Name),
		Role:       strings.TrimSpace(q.Role),
		DivisionID: q.DivisionID,
		IsActive:   q.IsActive,
	}, nil
}

func ToUserResponse(u *UserWithDivision, baseURL string) *UserResponse {
	avatarURL := buildFullURL(u.AvatarURL, baseURL)

	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: avatarURL,
		Phone:     u.Phone,
		Role:      u.Role,
		Division: Division{
			ID:   u.DivisionID,
			Name: u.DivisionName,
		},
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserResponses(users []UserWithDivision, baseURL string) []UserResponse {
	results := make([]UserResponse, len(users))
	for i, user := range users {
		results[i] = *ToUserResponse(&user, baseURL)
	}
	return results
}

func buildFullURL(relativePath *string, baseURL string) *string {
	if relativePath == nil || *relativePath == "" {
		return nil
	}
	fullURL := baseURL + *relativePath
	return &fullURL
}
