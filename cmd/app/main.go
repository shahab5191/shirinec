package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"shirinec.com/config"
	"shirinec.com/internal/db"
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/routes"
)

func main() {
	config.Load()

    database, err := db.NewDatabase()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    userRepo := repositories.NewUserRepository(database.Pool)
    incomeCategoryRepo := repositories.NewIncomeCategoryRepository(database.Pool)
    expenseCategoryRepo := repositories.NewExpenseCategoryRepository(database.Pool)
    
    deps := handler.Dependencies{
        UserRepo: userRepo,
        IncomeCategoryRepo: incomeCategoryRepo,
        ExpenseCategoryRepo: expenseCategoryRepo,
    }

    ginEngine := gin.Default()
    router := routes.NewRouter(ginEngine, &deps)
    router.SetupRouter()

	port := strconv.Itoa(config.AppConfig.Port)
	log.Printf("Starting %s in %s mode on port %s", config.AppConfig.AppName, config.AppConfig.Env, port)

    if err := ginEngine.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
