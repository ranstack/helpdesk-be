package middleware

import (
	"helpdesk/internal/utils/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		requestID := c.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		response.SetRequestID(c, requestID)
		c.Response().Header().Set("X-Request-ID", requestID)

		return next(c)
	}
}
