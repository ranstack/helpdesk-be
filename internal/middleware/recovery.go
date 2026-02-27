package middleware

import (
	"fmt"
	"helpdesk/internal/utils/response"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v5"
)

func Recovery(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					stack := debug.Stack()
					logger.Error("panic recovered",
						"error", err,
						"stack", string(stack),
						"uri", c.Request().URL.Path,
						"method", c.Request().Method,
					)

					c.JSON(http.StatusInternalServerError, response.Response{
						Error: &response.ErrorInfo{
							Message: "Internal server error",
						},
						Meta: &response.Meta{
							Timestamp: time.Now().UTC().Format(time.RFC3339),
						},
					})
				}
			}()

			return next(c)
		}
	}
}
