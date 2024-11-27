package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupUserRouter() {
	userService := services.NewUserService(r.Deps.UserRepo)
	userHandler := handler.NewUserHandler(userService)

	flags := middlewares.AuthMiddleWareFlags{ShouldBeActive: true}

	r.GinEngine.POST("/user/new_password", middlewares.AuthMiddleWare(flags, r.db), userHandler.NewPassword)
	r.GinEngine.POST("/user/new_email", middlewares.AuthMiddleWare(flags, r.db), userHandler.NewEmail)
    r.GinEngine.POST("/user/new_email/verify", middlewares.AuthMiddleWare(flags, r.db), userHandler.NewEmailVerification)
}
