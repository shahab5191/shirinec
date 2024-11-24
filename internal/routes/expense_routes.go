package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupIncomeRouter() {
    expenseService := services.NewExpenseService(
		r.Deps.ExpenseCategoryRepo,
	)
	expenseCategoryHandler := handler.NewExpenseCategoryHandler(expenseService)

	r.GinEngine.GET("/expense/category", middlewares.AuthMiddleWare(), expenseCategoryHandler.List)
	r.GinEngine.GET("/expense/category/:id", middlewares.AuthMiddleWare(), expenseCategoryHandler.GetByID)
    r.GinEngine.POST("/expense/category", middlewares.AuthMiddleWare(), expenseCategoryHandler.Create)
    r.GinEngine.DELETE("/expense/category/:id", middlewares.AuthMiddleWare(), expenseCategoryHandler.Delete)
    r.GinEngine.PUT("/expense/category/:id", middlewares.AuthMiddleWare(), expenseCategoryHandler.Update)
}
