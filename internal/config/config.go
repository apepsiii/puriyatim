package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	Env        string
	DBPath     string
	OneSender  OneSenderConfig
	Pakasir    PakasirConfig
	JWTSecret  string
	AppName    string
	AppVersion string
	AppURL     string
}

type OneSenderConfig struct {
	APIURL  string
	APIKey  string
	GroupID string
}

type PakasirConfig struct {
	ProjectSlug string
	APIKey      string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		DBPath:     getEnv("DB_PATH", "./db/puriyatim.db"),
		JWTSecret:  getEnv("JWT_SECRET", "default_secret_key"),
		AppName:    getEnv("APP_NAME", "Puri Yatim"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),
		AppURL:     getEnv("APP_URL", "http://localhost:8080"),
		OneSender: OneSenderConfig{
			APIURL:  getEnv("ONESENDER_API_URL", "https://api.onesender.com/v1/message"),
			APIKey:  getEnv("ONESENDER_API_KEY", ""),
			GroupID: getEnv("ONESENDER_GROUP_ID", ""),
		},
		Pakasir: PakasirConfig{
			ProjectSlug: getEnv("PAKASIR_PROJECT_SLUG", ""),
			APIKey:      getEnv("PAKASIR_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}