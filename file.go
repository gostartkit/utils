package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExist check file
func FileExist(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// DirExist check dir
func DirExist(dir string) bool {
	info, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// HomeDir join path with app home dir
func HomeDir(path string) (string, error) {
	if !filepath.IsAbs(path) {
		dir, err := os.Getwd()

		if err != nil {
			return "", err
		}

		path = filepath.Join(dir, path)
	}

	return path, nil
}

// FileSize format file size
func FileSize(size int64) string {

	units := []string{"B", "K", "M", "G", "T"}
	var i int
	var sf float64 = float64(size)

	for sf >= 1024 && i < len(units)-1 {
		sf /= 1024
		i++
	}

	return fmt.Sprintf("%.1f%s", sf, units[i])
}
