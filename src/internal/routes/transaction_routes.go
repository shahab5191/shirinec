package routes

import (
	handler "shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/middlewares"
	"shirinec.com/src/internal/services"
)

func (r *router) setupTransactionRouter() {
    transferService := services.NewTransferService(r.Deps.TransactionRepo)
    transferHandler := handler.NewTransferHandler(transferService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleWare := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/transfer", authMiddleWare, transferHandler.Transfer)
}
