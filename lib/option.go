package lib

import (
	"fmt"
	"os"
)

func ClearCache() error {
	cacheFile := GetCacheFilePath()
	err := os.Remove(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache tidak ditemukan")
		}
		return fmt.Errorf("gagal menghapus cache: %w", err)
	}
	return nil
}
