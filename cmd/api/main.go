package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"helpdesk/internal/config"
	"helpdesk/internal/database"
	"helpdesk/internal/features/category"
	"helpdesk/internal/features/division"
	"helpdesk/internal/middleware"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db := database.NewPostgres(cfg.DBConnString())
	defer db.Close()

	logger.Info("connected to database", "host", cfg.DBHost, "database", cfg.DBName)

	e := echo.New()

	e.Use(middleware.Recovery(logger))
	e.Use(middleware.Logger(logger))
	e.Use(middleware.CORS())

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo, logger)
	categoryHandler := category.NewHandler(categoryService)

	divisionRepo := division.NewRepository(db)
	divisionService := division.NewService(divisionRepo, logger)
	divisionHandler := division.NewHandler(divisionService)

	api := e.Group("/api/v1")

	api.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	category.RegisterRoutes(api, categoryHandler)
	division.RegisterRoutes(api, divisionHandler)
	addr := ":" + cfg.AppPort
	logger.Info("starting server", "address", addr, "app", cfg.AppName)
	fmt.Printf("🚀 Server started on %s\n", addr)

	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
