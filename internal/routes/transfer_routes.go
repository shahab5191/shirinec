package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupTransferRouter() {
    transferService := services.NewTransferService(r.Deps.TransactionRepo)
    transferHandler := handler.NewTransferHandler(transferService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleWare := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/transfer/:from/:to", authMiddleWare, transferHandler.Transfer)
}
