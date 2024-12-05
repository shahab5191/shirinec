package routes

import (
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupItemRouter() {
	itemService := services.NewItemService(&r.Deps.ItemRepo)
	itemHandler := handler.NewItemHandler(&itemService)

	flags := middlewares.AuthMiddleWareFlags{
		ShouldBeActive: true,
	}

	authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

	r.GinEngine.POST("/item", authMiddleware, itemHandler.Create)
	r.GinEngine.GET("/item", authMiddleware, itemHandler.List)
	r.GinEngine.GET("/item/:id", authMiddleware, itemHandler.GetByID)
	r.GinEngine.PUT("/item/:id", authMiddleware, itemHandler.Update)
	r.GinEngine.DELETE("/item/:id", authMiddleware, itemHandler.Delete)
}
