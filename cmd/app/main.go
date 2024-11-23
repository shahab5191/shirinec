package main

import (
	"log"
	"strconv"

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
    
    deps := handler.Dependencies{
        UserRepo: userRepo,
        IncomeCategoryRepo: incomeCategoryRepo,
    }

    r := routes.SetupRouter(deps)

	port := strconv.Itoa(config.AppConfig.Port)
	log.Printf("Starting %s in %s mode on port %s", config.AppConfig.AppName, config.AppConfig.Env, port)

    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
