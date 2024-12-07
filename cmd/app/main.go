package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"shirinec.com/config"
	"shirinec.com/internal/db"
	"shirinec.com/internal/handlers"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/routes"
	"shirinec.com/internal/utils"
	"shirinec.com/internal/validators"
	"shirinec.com/internal/workers"
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
	accountRepo := repositories.NewAccountRepository(database.Pool)
	mediaRepo := repositories.NewMediaRepository(database.Pool)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterValidators(v)
	}

	deps := handler.Dependencies{
		UserRepo:     userRepo,
		CategoryRepo: categoryRepo,
		ItemRepo:     itemRepo,
		AccountRepo:  accountRepo,
		MediaRepo:    mediaRepo,
	}

	utils.InitLogger()

	ginEngine := gin.Default()

	router := routes.NewRouter(ginEngine, &deps, database.Pool)
	router.SetupRouter()

	workers.ScheduleWorkers(mediaRepo)

	for _, route := range ginEngine.Routes() {
		log.Println(route.Method, route.Path)
	}
	port := strconv.Itoa(config.AppConfig.Port)
	log.Printf("Starting %s in %s mode on port %s", config.AppConfig.AppName, config.AppConfig.Env, port)

	if err := ginEngine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
