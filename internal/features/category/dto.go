package category

import (
	"helpdesk/internal/utils/validator"
	"time"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (r *CreateCategoryRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 20)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

func (r *UpdateCategoryRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 20)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCategoryResponse(c *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
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
