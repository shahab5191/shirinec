package routes

import (
	"google.golang.org/protobuf/internal/flags"
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
)

func (r *router) setupTransferRouter() {
    transferService := services.NewTransferSevice(r.Deps.TransferRepo)
    transferHandler := handler.NewTransferHandler(transferService)

    flags := middlewares.AuthMiddleWareFlags{
        ShouldBeActive: true,
    }

    authMiddleWare := middlewares.AuthMiddleWare(flags, r.db)

    r.GinEngine.POST("/transfer/:from/:to", authMiddleWare, transferHandler.Transfer)
}
