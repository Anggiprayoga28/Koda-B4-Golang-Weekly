package lib

import (
	"fmt"
	"os"
)

const cacheFile = "/tmp/menu_cache.json"

func ClearCache() error {
	err := os.Remove(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache tidak ditemukan")
		}
		return fmt.Errorf("gagal menghapus cache: %w", err)
	}
	return nil
}
