package routes

import (
	"log"

	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupItemRouter() {
    log.Println("Setting up item router")
	itemService := services.NewItemService(&r.Deps.ItemRepo)
	itemHandler := handler.NewItemHandler(&itemService)

	flags := middlewares.AuthMiddleWareFlags{
		ShouldBeActive: true,
	}

    r.GinEngine.POST("/item", middlewares.AuthMiddleWare(flags, r.db), itemHandler.Create)
    r.GinEngine.GET("/item", middlewares.AuthMiddleWare(flags, r.db), itemHandler.List)
}
