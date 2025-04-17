package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	// Load .env file
	godotenv.Load()

	Database()
	app()
	OAuth()
	// Set up Viper to automatically use environment variables
	viper.AutomaticEnv()
}

func app() {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")

	conf("TOKEN_SECRET_USER", "secret")
	conf("TOKEN_DURATION_USER", 24*time.Hour)
	conf("EMAIL_HOST", "")
	conf("EMAIL_PORT", "")
	conf("EMAIL_USERNAME", "")
	conf("EMAIL_PASSWORD", "")
}
