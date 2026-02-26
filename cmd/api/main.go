package main

import (
	"fmt"
	"helpdesk/internal/config"
	"helpdesk/internal/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.Load()
	db := database.NewPostgres(cfg.DBConnString())

	e := echo.New()

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	addr := ":" + cfg.AppPort
	fmt.Println("Server running on", addr)

	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}

	_ = db
}
