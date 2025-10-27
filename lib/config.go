package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	CacheDuration time.Duration
	CacheFilePath string
	APIURL        string
}

var AppConfig = &Config{
	CacheDuration: 15 * time.Minute,
	CacheFilePath: "/tmp/menu_cache.json",
	APIURL:        "https://raw.githubusercontent.com/Anggiprayoga28/Koda-B4-Golang--Weekly-Data/refs/heads/main/dataProduct.json",
}

func LoadConfig() error {
	file, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File .env tidak ditemukan, menggunakan konfigurasi default")
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

		switch key {
		case "CACHE_DURATION":
			if seconds, err := strconv.Atoi(value); err == nil {
				AppConfig.CacheDuration = time.Duration(seconds) * time.Second
			}
		case "CACHE_FILE_PATH":
			AppConfig.CacheFilePath = value
		case "API_URL":
			AppConfig.APIURL = value
		}
	}

	return scanner.Err()
}

func GetCacheDuration() time.Duration {
	return AppConfig.CacheDuration
}

func GetCacheFilePath() string {
	return AppConfig.CacheFilePath
}

func GetAPIURL() string {
	return AppConfig.APIURL
}
