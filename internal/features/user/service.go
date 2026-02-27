package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"helpdesk/internal/features/division"
	appErrors "helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/uploads"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAll(ctx context.Context, req *GetUsersQuery) (*response.ListResponse[UserResponse], error)
	GetByID(ctx context.Context, id int) (*UserResponse, error)
	Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, id int, req *UpdateUserRequest) (*UserResponse, error)
	UpdateAvatar(ctx context.Context, id int, avatarURL string) (*UserResponse, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo            Repository
	divisionService division.Service
	logger          *slog.Logger
	baseURL         string
}

func NewService(repo Repository, divisionService division.Service, logger *slog.Logger, baseURL string) Service {
	return &service{
		repo:            repo,
		divisionService: divisionService,
		logger:          logger,
		baseURL:         baseURL,
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
		Items: ToUserResponses(users, s.baseURL),
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

	return ToUserResponse(user, s.baseURL), nil
}

func (s *service) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)

	if err := s.divisionService.ValidateForAssignment(ctx, req.DivisionID); err != nil {
		return nil, err
	}

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

	role := strings.TrimSpace(req.Role)

	user, err := s.repo.Create(ctx, name, email, passwordHash, "", "", role, req.DivisionID)
	if err != nil {
		s.logger.Error("failed to create user", "error", err, "email", email)
		if strings.Contains(err.Error(), "already exists") {
			return nil, appErrors.AlreadyExists("User with this email")
		}
		return nil, appErrors.Internal("Failed to create user")
	}

	s.logger.Info("user created", "id", user.ID, "email", user.Email)
	return ToUserResponse(user, s.baseURL), nil
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

	if err := s.divisionService.ValidateForAssignment(ctx, req.DivisionID); err != nil {
		return nil, err
	}

	currentUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get current user", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update user")
	}
	if currentUser == nil {
		return nil, appErrors.NotFound("User")
	}

	name := strings.TrimSpace(req.Name)

	phone := ""
	if req.Phone != nil {
		phone = strings.TrimSpace(*req.Phone)
	}

	role := strings.TrimSpace(req.Role)

	isActive := currentUser.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	user, err := s.repo.Update(ctx, id, name, phone, role, req.DivisionID, isActive)
	if err != nil {
		s.logger.Error("failed to update user", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update user")
	}

	if user == nil {
		return nil, appErrors.NotFound("User")
	}

	s.logger.Info("user updated", "id", user.ID, "email", user.Email)
	return ToUserResponse(user, s.baseURL), nil
}

func (s *service) UpdateAvatar(ctx context.Context, id int, avatarURL string) (*UserResponse, error) {
	if id <= 0 {
		return nil, appErrors.BadRequest("Invalid user ID")
	}

	if avatarURL == "" {
		return nil, appErrors.BadRequest("Avatar URL is required")
	}

	oldUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update avatar")
	}
	if oldUser == nil {
		return nil, appErrors.NotFound("User")
	}

	user, err := s.repo.UpdateAvatar(ctx, id, avatarURL)
	if err != nil {
		s.logger.Error("failed to update avatar", "error", err, "id", id)
		return nil, appErrors.Internal("Failed to update avatar")
	}

	if user == nil {
		return nil, appErrors.NotFound("User")
	}

	if oldUser.AvatarURL != nil && *oldUser.AvatarURL != "" {
		if err := uploads.DeleteFile(*oldUser.AvatarURL); err != nil {
			s.logger.Warn("failed to delete old avatar", "error", err, "path", *oldUser.AvatarURL)
		}
	}

	s.logger.Info("user avatar updated", "id", user.ID)
	return ToUserResponse(user, s.baseURL), nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.BadRequest("Invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", "error", err, "id", id)
		return appErrors.Internal("Failed to delete user")
	}
	if user == nil {
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

	if user.AvatarURL != nil && *user.AvatarURL != "" {
		if err := uploads.DeleteFile(*user.AvatarURL); err != nil {
			s.logger.Warn("failed to delete user avatar", "error", err, "path", *user.AvatarURL)
		}
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
