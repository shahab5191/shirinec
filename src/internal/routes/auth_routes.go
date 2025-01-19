package routes

import (
	"shirinec.com/config"
	handler "shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/services"
)

func (r *router) setupAuthRouter() {
	authService := services.NewAuthService(
		r.Deps.UserRepo,
		config.AppConfig.JWTSecret,
	)
	authHandler := handler.NewAuthHandler(authService)

	r.GinEngine.POST("/auth/signup", authHandler.SignUp)
	r.GinEngine.POST("/auth/login", authHandler.Login)
	r.GinEngine.POST("/auth/refresh", authHandler.RefreshToken)
}
