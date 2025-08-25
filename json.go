package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ReadJSON read json to data
func ReadJSON(filename string, v any) error {

	f, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

// WriteJSON write data to json
func WriteJSON(filename string, v any, overwrite bool) error {

	if !overwrite && FileExist(filename) {
		return os.ErrExist
	}

	dir := filepath.Dir(filename)

	if !DirExist(dir) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)

	enc.SetIndent("", "  ")

	return enc.Encode(v)
}
