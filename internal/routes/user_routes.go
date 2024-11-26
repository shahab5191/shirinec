package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupUserRouter() {
	userService := services.NewUserService(r.Deps.UserRepo)
    userHandler := handler.NewUserHandler(userService)

    r.GinEngine.POST("/user/new_password", middlewares.AuthMiddleWare(), userHandler.NewPassword)
}
