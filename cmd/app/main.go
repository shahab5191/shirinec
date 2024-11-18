package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"shirinec.com/config"
)

func main() {
	config.Load()

	r := gin.Default()

	port := config.AppConfig.Port
	log.Printf("Starting %s in %s mode on port %s", config.AppConfig.AppName, config.AppConfig.Env, port)

    if err := r.Run(":" + port); err != nil {
    log.Fatalf("Failed to start server: %v", err)
    }
}
