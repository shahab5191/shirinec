package routes

import (
	"shirinec.com/config"
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/services"
)

func (r *router) setupAuthRouter() {
	authService := services.NewAuthService(
		r.Deps.UserRepo,
		config.AppConfig.JWTSecret,
	)
	authHandler := handler.NewAuthHandler(authService, r.validatorObj)

	r.GinEngine.POST("/auth/signup", authHandler.SignUp)
	r.GinEngine.POST("/auth/login", authHandler.Login)
	r.GinEngine.POST("/auth/refresh", authHandler.RefreshToken)
}
