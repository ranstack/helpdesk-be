package division

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

func (h *Handler) GetAll(c *echo.Context) error {
	var req GetDivisionsQuery
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.BadRequest("Invalid query parameters"))
	}

	divisions, err := h.service.GetAll(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Divisions retrieved successfully", divisions)
}

func (h *Handler) GetByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid division ID"))
	}

	division, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Division retrieved successfully", division)
}

func (h *Handler) Create(c *echo.Context) error {
	var req CreateDivisionRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	division, err := h.service.Create(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, "Division created successfully", division)
}

func (h *Handler) Update(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid division ID"))
	}

	var req UpdateDivisionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, err)
	}

	division, err := h.service.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.OK(c, "Division updated successfully", division)
}

func (h *Handler) Delete(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.Error(c, errors.BadRequest("Invalid division ID"))
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Division deleted successfully",
	})
}
