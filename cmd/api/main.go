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
	"helpdesk/internal/features/user"
	"helpdesk/internal/middleware"
	"helpdesk/internal/utils/uploads"

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

	if err := uploads.EnsureUploadDirs(); err != nil {
		log.Fatalf("failed to create upload directories: %v", err)
	}
	logger.Info("upload directories ready")

	e := echo.New()

	e.Use(middleware.RequestID)
	e.Use(middleware.Recovery(logger))
	e.Use(middleware.Logger(logger))
	e.Use(middleware.CORS())

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo, logger)
	categoryHandler := category.NewHandler(categoryService)

	divisionRepo := division.NewRepository(db)
	divisionService := division.NewService(divisionRepo, logger)
	divisionHandler := division.NewHandler(divisionService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, divisionService, logger, cfg.BaseURL)
	userHandler := user.NewHandler(userService)

	e.Static("/uploads", "uploads")

	api := e.Group("/api/v1")

	api.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	category.RegisterRoutes(api, categoryHandler)
	division.RegisterRoutes(api, divisionHandler)
	user.RegisterRoutes(api, userHandler)
	addr := ":" + cfg.AppPort
	logger.Info("starting server", "address", addr, "app", cfg.AppName)
	fmt.Printf("🚀 Server started on %s\n", addr)

	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
