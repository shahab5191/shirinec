package routes

import (
	"shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/middlewares"
	"shirinec.com/src/internal/services"
)

func (r *router) setupCategoryRouter() {
	categoryService := services.NewCategoryService(
		r.Deps.CategoryRepo,
	)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	flags := middlewares.AuthMiddleWareFlags{ShouldBeActive: true}

	r.GinEngine.GET("/category", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.List)
	r.GinEngine.GET("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.GetByID)
	r.GinEngine.POST("/category", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Create)
	r.GinEngine.DELETE("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Delete)
	r.GinEngine.PUT("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Update)
}
