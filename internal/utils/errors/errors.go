package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	CODE_NOT_FOUND        = "NOT_FOUND"
	CODE_ALREADY_EXISTS   = "ALREADY_EXISTS"
	CODE_VALIDATION_ERROR = "VALIDATION_ERROR"
	CODE_INTERNAL_ERROR   = "INTERNAL_SERVER_ERROR"
	CODE_BAD_REQUEST      = "BAD_REQUEST"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrValidation    = errors.New("validation error")
	ErrInternal      = errors.New("internal server error")
	ErrBadRequest    = errors.New("bad request")
)

type AppError struct {
	Err        error
	Code       string
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, code string, message string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func NotFound(resource string) *AppError {
	return &AppError{
		Err:        ErrNotFound,
		Code:       CODE_NOT_FOUND,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func AlreadyExists(resource string) *AppError {
	return &AppError{
		Err:        ErrAlreadyExists,
		Code:       CODE_ALREADY_EXISTS,
		Message:    fmt.Sprintf("%s already exists", resource),
		StatusCode: http.StatusConflict,
	}
}

func Validation(message string) *AppError {
	return &AppError{
		Err:        ErrValidation,
		Code:       CODE_VALIDATION_ERROR,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Err:        ErrInternal,
		Code:       CODE_INTERNAL_ERROR,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Err:        ErrBadRequest,
		Code:       CODE_BAD_REQUEST,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
