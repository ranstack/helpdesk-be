package division

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
	GetAll(ctx context.Context, req *GetDivisionsQuery) (*response.ListResponse[DivisionResponse], error)
	GetByID(ctx context.Context, id int) (*DivisionResponse, error)
	Create(ctx context.Context, req *CreateDivisionRequest) (*DivisionResponse, error)
	Update(ctx context.Context, id int, req *UpdateDivisionRequest) (*DivisionResponse, error)
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

func (s *service) GetAll(ctx context.Context, req *GetDivisionsQuery) (*response.ListResponse[DivisionResponse], error) {
	if req == nil {
		req = &GetDivisionsQuery{}
	}

	filter, err := req.Normalize()
	if err != nil {
		return nil, err
	}

	divisions, totalItems, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get divisions", "error", err)
		return nil, appErrors.Internal("Failed to retrieve divisions")
	}

	return &response.ListResponse[DivisionResponse]{
		Items: ToDivisionResponses(divisions),
		Pagination: response.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalItems: totalItems,
			TotalPages: response.CalculateTotalPages(totalItems, filter.Limit),
		},
	}, nil
}

func (s *service) GetByID(ctx context.Context, id int) (*DivisionResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid division ID")
	}

	division, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get division", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to retrieve division")
	}

	if division == nil {
		return nil, appErrors.NotFound("Division")
	}

	return ToDivisionResponse(division), nil
}

func (s *service) Create(ctx context.Context, req *CreateDivisionRequest) (*DivisionResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("failed to check existing division", "error", err)
		return nil, appErrors.Internal("Failed to create division")
	}
	if existing != nil {
		return nil, appErrors.AlreadyExists("Division")
	}

	division, err := s.repo.Create(ctx, name)
	if err != nil {
		s.logger.Error("failed to create division", "error", err, "name", name)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("Division")
		}
		return nil, appErrors.Internal("Failed to create division")
	}

	s.logger.Info("division created", "id", division.ID, "name", division.Name)
	return ToDivisionResponse(division), nil
}

func (s *service) Update(ctx context.Context, id int, req *UpdateDivisionRequest) (*DivisionResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid division ID")
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check division existence", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update division")
	}
	if !exists {
		return nil, appErrors.NotFound("Division")
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("failed to check existing division", "error", err)
		return nil, appErrors.Internal("Failed to update division")
	}
	if existing != nil && existing.ID != id {
		return nil, appErrors.AlreadyExists("Division with this name")
	}

	division, err := s.repo.Update(ctx, id, name)
	if err != nil {
		s.logger.Error("failed to update division", "error", err, "id", id)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("Division")
		}
		return nil, appErrors.Internal("Failed to update division")
	}

	if division == nil {
		return nil, appErrors.NotFound("Division")
	}

	s.logger.Info("division updated", "id", division.ID, "name", division.Name)
	return ToDivisionResponse(division), nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.BadRequest("Invalid division ID")
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check division existence", "error", err, "id", id)
		return appErrors.Internal("Failed to delete division")
	}
	if !exists {
		return appErrors.NotFound("Division")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.NotFound("Division")
		}
		s.logger.Error("failed to delete division", "error", err, "id", id)
		return appErrors.Internal(fmt.Sprintf("Failed to delete division: %v", err))
	}

	s.logger.Info("division deleted", "id", id)
	return nil
}
