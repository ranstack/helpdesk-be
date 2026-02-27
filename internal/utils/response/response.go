package response

import (
	"helpdesk/internal/utils/errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func Success(c *echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusCreated, message, data)
}

func OK(c *echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusOK, message, data)
}

func NoContent(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Error(c *echo.Context, err error) error {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.Internal(err.Error())
	}

	errorInfo := &ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	}

	return c.JSON(appErr.StatusCode, Response{
		Success: false,
		Error:   errorInfo,
	})
}
