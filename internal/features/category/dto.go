package category

import (
	appErrors "helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/validator"
	"strings"
	"time"
)

const (
	defaultCategoryPage  = 1
	defaultCategoryLimit = 10
	maxCategoryLimit     = 100
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCategoriesQuery struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Name      string `query:"name"`
	IsActive  *bool  `query:"isActive"`
	CreatedAt string `query:"createdAt"`
}

type CategoryListFilter struct {
	Page      int
	Limit     int
	Offset    int
	Name      string
	IsActive  *bool
	CreatedAt *time.Time
}

type CategoryListResponse struct {
	Items      []CategoryResponse          `json:"items"`
	Pagination response.PaginationResponse `json:"pagination"`
}

func (r *CreateCategoryRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 20)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (r *UpdateCategoryRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 20)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (q *GetCategoriesQuery) Normalize() (*CategoryListFilter, error) {
	page := q.Page
	if page == 0 {
		page = defaultCategoryPage
	}
	if page < 1 {
		return nil, appErrors.BadRequest("page must be greater than 0")
	}

	limit := q.Limit
	if limit == 0 {
		limit = defaultCategoryLimit
	}
	if limit < 1 {
		return nil, appErrors.BadRequest("limit must be greater than 0")
	}
	if limit > maxCategoryLimit {
		limit = maxCategoryLimit
	}

	var createdAt *time.Time
	if strings.TrimSpace(q.CreatedAt) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(q.CreatedAt))
		if err != nil {
			return nil, appErrors.BadRequest("createdAt must use YYYY-MM-DD format")
		}
		createdAt = &parsed
	}

	return &CategoryListFilter{
		Page:      page,
		Limit:     limit,
		Offset:    (page - 1) * limit,
		Name:      strings.TrimSpace(q.Name),
		IsActive:  q.IsActive,
		CreatedAt: createdAt,
	}, nil
}

func ToCategoryResponse(c *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
	}
}

func ToCategoryResponses(categories []Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = *ToCategoryResponse(&c)
	}
	return responses
}
