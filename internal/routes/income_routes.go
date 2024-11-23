package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupIncomeRouter() {
    incomeService := services.NewIncomeService(
		r.Deps.IncomeCategoryRepo,
	)
	incomeCategoryHandler := handler.NewIncomeCategoryHandler(incomeService)

	r.GinEngine.GET("/income/category", middlewares.AuthMiddleWare(), incomeCategoryHandler.List)
	r.GinEngine.GET("/income/category/:id", middlewares.AuthMiddleWare(), incomeCategoryHandler.GetByID)
    r.GinEngine.POST("/income/category", middlewares.AuthMiddleWare(), incomeCategoryHandler.Create)
    r.GinEngine.DELETE("/income/category/:id", middlewares.AuthMiddleWare(), incomeCategoryHandler.Delete)
    r.GinEngine.PUT("/income/category/:id", middlewares.AuthMiddleWare(), incomeCategoryHandler.Update)
}
