package routes

import (
	"shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/middlewares"
	"shirinec.com/src/internal/services"
)

func (r *router) setupAccountRouter() {
    accountService := services.NewAccountService(&r.Deps.AccountRepo)
    accountHandler := handler.NewAccountHandler(&accountService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/account", authMiddleware, accountHandler.Create)
	r.GinEngine.GET("/account", authMiddleware, accountHandler.List)
	r.GinEngine.GET("/account/:id", authMiddleware, accountHandler.GetByID)
    r.GinEngine.PUT("/account/:id", authMiddleware, accountHandler.Update)
    r.GinEngine.DELETE("/account/:id", authMiddleware, accountHandler.Delete)

}
