package lib

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

var AppConfig = &Config{}

func LoadConfig() error {
	godotenv.Load()

	AppConfig.DatabaseURL = os.Getenv("DATABASE_URL")

	if AppConfig.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL tidak ditemukan")
	}

	return nil
}

func GetDatabaseURL() string {
	return AppConfig.DatabaseURL
}
