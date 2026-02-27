package category

import (
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/validator"
	"strings"
	"time"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"isActive"`
}

type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCategoriesQuery struct {
	response.PaginationQuery
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
	page, limit, offset := q.NormalizePagination()

	createdAt, err := response.ParseDate(q.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &CategoryListFilter{
		Page:      page,
		Limit:     limit,
		Offset:    offset,
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
	return response.MapResponses(categories, ToCategoryResponse)
}
