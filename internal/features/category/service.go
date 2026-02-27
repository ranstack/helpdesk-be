package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	appErrors "helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"
)

type Service interface {
	GetAll(ctx context.Context, req *GetCategoriesQuery) (*response.ListResponse[CategoryResponse], error)
	GetByID(ctx context.Context, id int) (*CategoryResponse, error)
	Create(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error)
	Update(ctx context.Context, id int, req *UpdateCategoryRequest) (*CategoryResponse, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) GetAll(ctx context.Context, req *GetCategoriesQuery) (*response.ListResponse[CategoryResponse], error) {
	if req == nil {
		req = &GetCategoriesQuery{}
	}

	filter, err := req.Normalize()
	if err != nil {
		return nil, err
	}

	categories, totalItems, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get categories", "error", err)
		return nil, appErrors.Internal("Failed to retrieve categories")
	}

	return &response.ListResponse[CategoryResponse]{
		Items: ToCategoryResponses(categories),
		Pagination: response.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalItems: totalItems,
			TotalPages: response.CalculateTotalPages(totalItems, filter.Limit),
		},
	}, nil
}

func (s *service) GetByID(ctx context.Context, id int) (*CategoryResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid category ID")
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get category", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to retrieve category")
	}

	if category == nil {
		return nil, appErrors.NotFound("Category")
	}

	return ToCategoryResponse(category), nil
}

func (s *service) Create(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("failed to check existing category", "error", err)
		return nil, appErrors.Internal("Failed to create category")
	}
	if existing != nil {
		return nil, appErrors.AlreadyExists("Category")
	}

	category, err := s.repo.Create(ctx, name)
	if err != nil {
		s.logger.Error("failed to create category", "error", err, "name", name)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("Category")
		}
		return nil, appErrors.Internal("Failed to create category")
	}

	s.logger.Info("category created", "id", category.ID, "name", category.Name)
	return ToCategoryResponse(category), nil
}

func (s *service) Update(ctx context.Context, id int, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid category ID")
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check category existence", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update category")
	}
	if !exists {
		return nil, appErrors.NotFound("Category")
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("failed to check existing category", "error", err)
		return nil, appErrors.Internal("Failed to update category")
	}
	if existing != nil && existing.ID != id {
		return nil, appErrors.AlreadyExists("Category with this name")
	}

	category, err := s.repo.Update(ctx, id, name)
	if err != nil {
		s.logger.Error("failed to update category", "error", err, "id", id)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("Category")
		}
		return nil, appErrors.Internal("Failed to update category")
	}

	if category == nil {
		return nil, appErrors.NotFound("Category")
	}

	s.logger.Info("category updated", "id", category.ID, "name", category.Name)
	return ToCategoryResponse(category), nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.BadRequest("Invalid category ID")
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check category existence", "error", err, "id", id)
		return appErrors.Internal("Failed to delete category")
	}
	if !exists {
		return appErrors.NotFound("Category")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.NotFound("Category")
		}
		s.logger.Error("failed to delete category", "error", err, "id", id)
		return appErrors.Internal(fmt.Sprintf("Failed to delete category: %v", err))
	}

	s.logger.Info("category deleted", "id", id)
	return nil
}
