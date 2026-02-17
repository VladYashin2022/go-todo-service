package storage

import (
	"errors"
	"fmt"
	"os"
)

const storagePath = "storage.json"

func WriteFile() (*os.File, error) {
	file, err := os.OpenFile(storagePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file to write: %w", err)
	}
	return file, nil
}

func ReadFile() ([]byte, error) {
	data, err := os.ReadFile(storagePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return data, nil
}
