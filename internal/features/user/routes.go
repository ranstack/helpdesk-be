package user

import "github.com/labstack/echo/v5"

func RegisterRoutes(g *echo.Group, handler *Handler) {
	users := g.Group("/users")

	users.GET("", handler.GetAll)
	users.GET("/:id", handler.GetByID)
	users.POST("", handler.Create)
	users.PATCH("/:id", handler.Update)
	users.DELETE("/:id", handler.Delete)
}
