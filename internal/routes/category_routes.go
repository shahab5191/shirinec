package routes

import (
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupCategoryRouter() {
	categoryService := services.NewCategoryService(
		r.Deps.CategoryRepo,
	)
	categoryHandler := handler.NewCategoryHandler(categoryService, r.validatorObj)

	flags := middlewares.AuthMiddleWareFlags{ShouldBeActive: true}

	r.GinEngine.GET("/category", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.List)
	r.GinEngine.GET("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.GetByID)
	r.GinEngine.POST("/category", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Create)
	r.GinEngine.DELETE("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Delete)
	r.GinEngine.PUT("/category/:id", middlewares.AuthMiddleWare(flags, r.db), categoryHandler.Update)
}
