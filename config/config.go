package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppName              string
	Port                 int
	Env                  string
	DatabaseURL          string
	JWTSecret            string
	JWTRefreshSecret     string
	PoolSize             int
	Timeout              time.Duration
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	RedisURL             string
	UploadFolder         string
	SqlFolder            string
}

var AppConfig *Config

func Load() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found!")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetDefault("AppName", "PersonalFinanceApp")
	viper.SetDefault("Port", 5500)
	viper.SetDefault("Env", "development")
	viper.SetDefault("PoolSize", 10)
	viper.SetDefault("Timeout", 5*time.Second)
	viper.SetDefault("AccessTokenRefresh", 15*time.Minute)
	viper.SetDefault("RefreshTokenDuration", 168*time.Hour)
	viper.SetDefault("UploadFolder", "./upload")
	viper.SetDefault("SqlFolder", "./internal/db/sql")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("No config file found: %s", err)
	}

	AppConfig = &Config{
		AppName:              viper.GetString("AppName"),
		Port:                 viper.GetInt("server.port"),
		Env:                  viper.GetString("server.env"),
		PoolSize:             viper.GetInt("database.pool_size"),
		Timeout:              viper.GetDuration("database.timeout"),
		DatabaseURL:          getEnvOrDefault("DATABASE_URL", ""),
		JWTSecret:            getEnvOrDefault("JWT_SECRET", ""),
		JWTRefreshSecret:     getEnvOrDefault("JWT_REFRESH_SECRET", ""),
		AccessTokenDuration:  viper.GetDuration("services.auth.access_token_duration"),
		RefreshTokenDuration: viper.GetDuration("services.auth.refresh_token_duration"),
		RedisURL:             getEnvOrDefault("REDIS_URL", ""),
		UploadFolder:         viper.GetString("server.upload_folder"),
		SqlFolder:            viper.GetString("server.sql_folder"),
	}
	println(viper.GetInt("database.pool_size"))

	fmt.Printf("%+v", AppConfig)

	if AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	if AppConfig.RedisURL == "" {
		log.Fatal("REDIS_URL is required but not set")
	}

	log.Printf("Config loaded: %+v", AppConfig)
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
