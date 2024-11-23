package routes

import (
	"shirinec.com/config"
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/services"
)

func (r *router) setupAuthRouter() {
	userService := services.NewAuthService(
		r.Deps.UserRepo,
		config.AppConfig.JWTSecret,
	)
	authHandler := handler.NewAuthHandler(userService)

	r.GinEngine.POST("/auth/signup", authHandler.SignUp)
	r.GinEngine.POST("/auth/login", authHandler.Login)
	r.GinEngine.POST("/auth/refresh", authHandler.RefreshToken)
}
