package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIURL      string
	DatabaseURL string
}

var AppConfig = &Config{
	DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		AppConfig.DatabaseURL = dbURL
		fmt.Println("Menggunakan DATABASE_URL dari environment variable")
	}

	if apiURL := os.Getenv("API_URL"); apiURL != "" {
		AppConfig.APIURL = apiURL
		fmt.Println("Menggunakan API_URL dari environment variable")
	}

	file, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File .env tidak ditemukan, menggunakan environment variables atau default")
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		updateConfig(&key, &value)
	}

	return scanner.Err()
}

func updateConfig(key *string, value *string) {
	switch *key {
	case "API_URL":
		AppConfig.APIURL = *value
	case "DATABASE_URL":
		AppConfig.DatabaseURL = *value
	}
}

func GetAPIURL() string {
	return AppConfig.APIURL
}

func GetDatabaseURL() string {
	return AppConfig.DatabaseURL
}
