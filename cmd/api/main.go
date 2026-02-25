package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.GET("/", hello)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
