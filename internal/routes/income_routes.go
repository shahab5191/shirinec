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


}
