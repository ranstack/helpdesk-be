package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func Logger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res, _ := echo.UnwrapResponse(c.Response())
			status := 0
			if res != nil {
				status = res.Status
			}

			logger.Info("request",
				"method", req.Method,
				"uri", req.URL.Path,
				"status", status,
				"latency", time.Since(start).String(),
				"ip", c.RealIP(),
				"user_agent", req.UserAgent(),
			)

			return err
		}
	}
}
