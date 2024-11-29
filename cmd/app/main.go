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

	db.NewRedis()

	userRepo := repositories.NewUserRepository(database.Pool)
	categoryRepo := repositories.NewCategoryRepository(database.Pool)
	itemRepo := repositories.NewItemRepository(database.Pool)

	deps := handler.Dependencies{
		UserRepo:     userRepo,
		CategoryRepo: categoryRepo,
		ItemRepo:     itemRepo,
	}

	ginEngine := gin.Default()
	router := routes.NewRouter(ginEngine, &deps, database.Pool)
	router.SetupRouter()

    for _, route := range ginEngine.Routes() {
        log.Println(route.Method, route.Path)
    }
	port := strconv.Itoa(config.AppConfig.Port)
	log.Printf("Starting %s in %s mode on port %s", config.AppConfig.AppName, config.AppConfig.Env, port)

	if err := ginEngine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
