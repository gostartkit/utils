package utils

import (
	"io"
	"net/http"
	"os"
	"time"
)

// FetchFile
func FetchFile(url string, dist string) error {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.Create(dist)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

// TryFetchFile
func TryFetchFile(url string, dist string, retry int) error {

	var err error

	for i := retry; i > 0; i-- {

		if err = FetchFile(url, dist); err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return err
	}

	return nil
}
