package category

import (
	"helpdesk/internal/utils/errors"
	"helpdesk/internal/utils/response"
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

func (h *Handler) Create(c *echo.Context) error {
	var req CreateCategoryRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	category, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "Category created successfully", category)
}

func (h *Handler) GetByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid category ID"))
	}

	category, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Category retrieved successfully", category)
}

func (h *Handler) GetAll(c *echo.Context) error {
	categories, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Categories retrieved successfully", categories)
}

func (h *Handler) Update(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid category ID"))
	}

	var req UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	category, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Category updated successfully", category)
}

func (h *Handler) Delete(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid category ID"))
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Category deleted successfully",
	})
}
