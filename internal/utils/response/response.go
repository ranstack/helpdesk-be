package response

import (
	"helpdesk/internal/utils/errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Meta struct {
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta"`
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

type ListResponse[T any] struct {
	Items      []T                `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *PaginationQuery) NormalizePagination() (page int, limit int, offset int) {
	page = p.Page
	if page == 0 {
		page = DefaultPage
	}
	if page < 1 {
		page = DefaultPage
	}

	limit = p.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	if limit < 1 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset = (page - 1) * limit
	return
}

func ParseDate(dateStr string) (*time.Time, error) {
	if strings.TrimSpace(dateStr) == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return nil, errors.BadRequest("Date must use YYYY-MM-DD format")
	}

	return &parsed, nil
}

func CalculateTotalPages(totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}
	return (totalItems + limit - 1) / limit
}

func GetRequestID(c *echo.Context) string {
	if c == nil {
		return uuid.New().String()
	}
	if id, ok := c.Get("requestId").(string); ok {
		return id
	}
	return uuid.New().String()
}

func SetRequestID(c *echo.Context, requestID string) {
	if c != nil {
		c.Set("requestId", requestID)
	}
}

func buildMeta(c *echo.Context) *Meta {
	return &Meta{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func Success(c *echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Message: message,
		Data:    data,
		Meta:    buildMeta(c),
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
		Error: errorInfo,
		Meta:  buildMeta(c),
	})
}

func MapResponses[T any, R any](items []T, mapper func(*T) *R) []R {
	responses := make([]R, len(items))
	for i, item := range items {
		responses[i] = *mapper(&item)
	}
	return responses
}
