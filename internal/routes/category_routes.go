package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupCategoryRouter() {
    categoryService := services.NewCategoryService(
		r.Deps.CategoryRepo,
	)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	r.GinEngine.GET("/category", middlewares.AuthMiddleWare(), categoryHandler.List)
	r.GinEngine.GET("/category/:id", middlewares.AuthMiddleWare(), categoryHandler.GetByID)
    r.GinEngine.POST("/category", middlewares.AuthMiddleWare(), categoryHandler.Create)
    r.GinEngine.DELETE("/category/:id", middlewares.AuthMiddleWare(), categoryHandler.Delete)
    r.GinEngine.PUT("/category/:id", middlewares.AuthMiddleWare(), categoryHandler.Update)
}
