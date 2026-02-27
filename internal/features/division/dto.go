package division

import (
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/validator"
	"strings"
	"time"
)

type CreateDivisionRequest struct {
	Name string `json:"name"`
}

type UpdateDivisionRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type DivisionResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetDivisionsQuery struct {
	response.PaginationQuery
	Name      string `query:"name"`
	IsActive  *bool  `query:"isActive"`
	CreatedAt string `query:"createdAt"`
}

type DivisionListFilter struct {
	Page      int
	Limit     int
	Offset    int
	Name      string
	IsActive  *bool
	CreatedAt *time.Time
}

func (r *CreateDivisionRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 50)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (r *UpdateDivisionRequest) Validate() error {
	v := validator.New()

	validator.ValidateString(v, "name", r.Name, true, 2, 50)

	if !v.Valid() {
		return v.ToAppError()
	}

	return nil
}

func (q *GetDivisionsQuery) Normalize() (*DivisionListFilter, error) {
	page, limit, offset := q.NormalizePagination()

	createdAt, err := response.ParseDate(q.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &DivisionListFilter{
		Page:      page,
		Limit:     limit,
		Offset:    offset,
		Name:      strings.TrimSpace(q.Name),
		IsActive:  q.IsActive,
		CreatedAt: createdAt,
	}, nil
}

func ToDivisionResponse(d *Division) *DivisionResponse {
	return &DivisionResponse{
		ID:        d.ID,
		Name:      d.Name,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
	}
}

func ToDivisionResponses(divisions []Division) []DivisionResponse {
	responses := make([]DivisionResponse, len(divisions))
	for i, d := range divisions {
		responses[i] = *ToDivisionResponse(&d)
	}
	return responses
}
