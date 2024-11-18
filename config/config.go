package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName      string
	Port         string
	Env          string
	DatabaseURL  string
	JWTSecret    string
}

var AppConfig *Config


func Load(){
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.SetDefault("AppName", "PersonalFinanceApp")
    viper.SetDefault("Port", 5500)
    viper.SetDefault("Env", "development")

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Fatal("No config file found: %v", err)
    }

    AppConfig = &Config{
        AppName:        viper.GetString("AppName"),
        Port:           viper.GetString("Port"),
        Env:            viper.GetString("Env"),
        DatabaseURL:    viper.GetString("DatabaseURL"),
        JWTSecret:      viper.GetString("JWTSecret"),
    }

    if AppConfig.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is required but not set")
    }

    if AppConfig.JWTSecret == "" {
        log.Fatal("JWT_SECRET is required but not set")
    }

    log.Printf("Config loaded: %+v", AppConfig)
}
