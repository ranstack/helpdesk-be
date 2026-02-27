package user

import (
	"helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"
	"helpdesk/internal/utils/uploads"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *echo.Context) error {
	var req GetUsersQuery
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid query parameters"))
	}

	users, err := h.service.GetAll(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Users retrieved successfully", users)
}

func (h *Handler) GetByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid user ID"))
	}

	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "User retrieved successfully", user)
}

func (h *Handler) Create(c *echo.Context) error {
	var req CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	user, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "User created successfully", user)
}

func (h *Handler) Update(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid user ID"))
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	user, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "User updated successfully", user)
}

func (h *Handler) UpdateAvatar(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid user ID"))
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.Error(c, errors.BadRequest("Avatar file is required"))
	}

	avatarURL, err := uploads.SaveAvatarImage(fileHeader)
	if err != nil {
		return response.Error(c, err)
	}

	user, err := h.service.UpdateAvatar(c.Request().Context(), id, avatarURL)
	if err != nil {
		uploads.DeleteFile(avatarURL)
		return response.Error(c, err)
	}

	return response.OK(c, "Avatar updated successfully", user)
}

func (h *Handler) Delete(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid user ID"))
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}
