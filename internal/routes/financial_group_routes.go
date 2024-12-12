package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupFinancialGroupRouter() {
    financialGroupService := services.NewFinancialGroupService(&r.Deps.FinancialGroupRepo)
    financialGroupHandler := handler.NewFinancialGroupHandler(&financialGroupService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/financial_group", authMiddleware, financialGroupHandler.Create)
    r.GinEngine.POST("/financial_group/:id/add_user", authMiddleware, financialGroupHandler.AddUser)
}
