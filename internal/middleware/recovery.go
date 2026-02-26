package middleware

import (
	"fmt"
	"helpdesk/internal/utils/response"
	"log/slog"
	"net/http"
	"runtime/debug"

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
						Success: false,
						Error: &response.ErrorInfo{
							Message: "Internal server error",
						},
					})
				}
			}()

			return next(c)
		}
	}
}
