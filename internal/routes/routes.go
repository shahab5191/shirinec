package routes

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/config"
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func SetupRouter(deps handler.Dependencies) *gin.Engine {
    r := gin.Default()
    userService := services.NewAuthService(
        deps.UserRepo,
        config.AppConfig.JWTSecret,
    )
    incomeService := services.NewIncomeService(
        deps.IncomeCategoryRepo,
    )
    authHandler := handler.NewAuthHandler(userService)
    incomeCategoryHandler := handler.NewIncomeCategoryHandler(incomeService)

    r.POST("/auth/signup", authHandler.SignUp)
    r.POST("/auth/login", authHandler.Login)
    r.POST("/auth/refresh", authHandler.RefreshToken)

    r.GET("/income/category", middlewares.AuthMiddleWare(), incomeCategoryHandler.List)
    return r
}
