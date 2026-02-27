package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	appErrors "helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAll(ctx context.Context, req *GetUsersQuery) (*response.ListResponse[UserResponse], error)
	GetByID(ctx context.Context, id int) (*UserResponse, error)
	Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, id int, req *UpdateUserRequest) (*UserResponse, error)
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

func (s *service) GetAll(ctx context.Context, req *GetUsersQuery) (*response.ListResponse[UserResponse], error) {
	if req == nil {
		req = &GetUsersQuery{}
	}

	filter, err := req.Normalize()
	if err != nil {
		return nil, err
	}

	users, totalItems, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get users", "error", err)
		return nil, appErrors.Internal("Failed to retrieve users")
	}

	return &response.ListResponse[UserResponse]{
		Items: ToUserResponses(users),
		Pagination: response.PaginationResponse{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalItems: totalItems,
			TotalPages: response.CalculateTotalPages(totalItems, filter.Limit),
		},
	}, nil
}

func (s *service) GetByID(ctx context.Context, id int) (*UserResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to retrieve user")
	}

	if user == nil {
		return nil, appErrors.NotFound("User")
	}

	return ToUserResponse(user), nil
}

func (s *service) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)

	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to check existing user", "error", err)
		return nil, appErrors.Internal("Failed to create user")
	}
	if existing != nil {
		return nil, appErrors.AlreadyExists("User with this email")
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, appErrors.Internal("Failed to create user")
	}

	avatarURL := ""
	if req.AvatarURL != nil {
		avatarURL = strings.TrimSpace(*req.AvatarURL)
	}

	phone := ""
	if req.Phone != nil {
		phone = strings.TrimSpace(*req.Phone)
	}

	role := strings.TrimSpace(req.Role)

	user, err := s.repo.Create(ctx, name, email, passwordHash, avatarURL, phone, role, req.DivisionID)
	if err != nil {
		s.logger.Error("failed to create user", "error", err, "email", email)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("User with this email")
		}
		return nil, appErrors.Internal("Failed to create user")
	}

	s.logger.Info("user created", "id", user.ID, "email", user.Email)
	return ToUserResponse(user), nil
}

func (s *service) Update(ctx context.Context, id int, req *UpdateUserRequest) (*UserResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid user ID")
	}

	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check user existence", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update user")
	}
	if !exists {
		return nil, appErrors.NotFound("User")
	}

	name := strings.TrimSpace(req.Name)

	avatarURL := ""
	if req.AvatarURL != nil {
		avatarURL = strings.TrimSpace(*req.AvatarURL)
	}

	phone := ""
	if req.Phone != nil {
		phone = strings.TrimSpace(*req.Phone)
	}

	role := strings.TrimSpace(req.Role)

	user, err := s.repo.Update(ctx, id, name, avatarURL, phone, role, req.DivisionID, req.IsActive)
	if err != nil {
		s.logger.Error("failed to update user", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update user")
	}

	if user == nil {
		return nil, appErrors.NotFound("User")
	}

	s.logger.Info("user updated", "id", user.ID, "email", user.Email)
	return ToUserResponse(user), nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.BadRequest("Invalid user ID")
	}

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check user existence", "error", err, "id", id)
		return appErrors.Internal("Failed to delete user")
	}
	if !exists {
		return appErrors.NotFound("User")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.NotFound("User")
		}
		s.logger.Error("failed to delete user", "error", err, "id", id)
		return appErrors.Internal(fmt.Sprintf("Failed to delete user: %v", err))
	}

	s.logger.Info("user deleted", "id", id)
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
