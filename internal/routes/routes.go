package routes

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/config"
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/services"
)

func SetupRouter(deps handler.Dependencies) *gin.Engine {
    r := gin.Default()
    userHandler := handler.NewUserHandler(deps.UserRepo)
    userService := services.NewAuthService(
        deps.UserRepo,
        config.AppConfig.JWTSecret,
    )
    authHandler := handler.NewAuthHandler(*userService)

    r.POST("/auth/signup", userHandler.CreateUser)
    r.POST("/auth/login", authHandler.Login)
    return r
}
