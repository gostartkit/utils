package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

// TarFile compress file with gzip
func TarFile(src string) (string, error) {

	sf, err := os.Open(src)

	if err != nil {
		return "", err
	}

	defer sf.Close()

	dist := src + ".tar.gz"

	df, err := os.Create(dist)

	if err != nil {
		return "", err
	}

	defer df.Close()

	gw := gzip.NewWriter(df)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	fi, err := sf.Stat()

	if err != nil {
		return "", err
	}

	header := &tar.Header{
		Name: fi.Name(),
		Size: fi.Size(),
		Mode: int64(fi.Mode()),
	}

	if err := tw.WriteHeader(header); err != nil {
		return "", err
	}

	if _, err := io.Copy(tw, sf); err != nil {
		return "", err
	}

	return dist, nil
}
