package routes

import (
	handler "shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/middlewares"
	"shirinec.com/src/internal/services"
)

func (r *router) setupFinancialGroupRouter() {
    financialGroupService := services.NewFinancialGroupService(&r.Deps.FinancialGroupRepo)
    financialGroupHandler := handler.NewFinancialGroupHandler(&financialGroupService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/financial_group", authMiddleware, financialGroupHandler.Create)
    r.GinEngine.POST("/financial_group/:id/:userID", authMiddleware, financialGroupHandler.AddUser)
    r.GinEngine.GET("/financial_group/:id", authMiddleware, financialGroupHandler.GetByID)
    r.GinEngine.GET("/financial_group", authMiddleware, financialGroupHandler.List)
    r.GinEngine.DELETE("/financial_group/:id", authMiddleware, financialGroupHandler.Delete)
    r.GinEngine.DELETE("/financial_group/:id/:userID", authMiddleware, financialGroupHandler.RemoveGroupMember)
}
