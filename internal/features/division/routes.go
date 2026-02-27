package division

import "github.com/labstack/echo/v5"

func RegisterRoutes(g *echo.Group, handler *Handler) {
	divisions := g.Group("/divisions")

	divisions.GET("", handler.GetAll)
	divisions.GET("/:id", handler.GetByID)
	divisions.POST("", handler.Create)
	divisions.PATCH("/:id", handler.Update)
	divisions.DELETE("/:id", handler.Delete)
}
