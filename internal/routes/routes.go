package routes

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/config"
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/services"
)

func SetupRouter(deps handler.Dependencies) *gin.Engine {
    r := gin.Default()
    userService := services.NewAuthService(
        deps.UserRepo,
        config.AppConfig.JWTSecret,
    )
    authHandler := handler.NewAuthHandler(*userService)

    r.POST("/auth/signup", authHandler.SignUp)
    r.POST("/auth/login", authHandler.Login)
    r.POST("/auth/refresh", authHandler.RefreshToken)
    return r
}
