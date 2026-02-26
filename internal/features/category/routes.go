package category

import "github.com/labstack/echo/v5"

func RegisterRoutes(g *echo.Group, handler *Handler) {
	categories := g.Group("/categories")

	categories.POST("", handler.Create)
	categories.GET("", handler.GetAll)
	categories.GET("/:id", handler.GetByID)
	categories.PUT("/:id", handler.Update)
	categories.DELETE("/:id", handler.Delete)
}
