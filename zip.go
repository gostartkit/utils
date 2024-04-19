package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

// ZipFile compress file with zip
func ZipFile(src string) (string, error) {

	srcFile, err := os.Open(src)

	if err != nil {
		return "", err
	}

	defer srcFile.Close()

	dist := src + ".zip"

	distFile, err := os.Create(dist)

	if err != nil {
		return "", err
	}

	defer distFile.Close()

	zipWriter := zip.NewWriter(distFile)
	defer zipWriter.Close()

	w, err := zipWriter.Create(path.Base(src))

	if err != nil {
		return "", err
	}

	if _, err := io.Copy(w, srcFile); err != nil {
		return "", err
	}

	return dist, nil
}
