package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupIncomeRouter() {
	incomeService := services.NewIncomeService(r.Deps.TransactionRepo)
	incomeHandler := handler.NewIncomeHandler(incomeService)

	flags := middlewares.AuthMiddleWareFlags{
		ShouldBeActive: true,
	}

	authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

	r.GinEngine.POST("/income", authMiddleware, incomeHandler.Create)
}
