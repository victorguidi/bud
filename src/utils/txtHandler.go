package utils

import (
	"io"
	"os"
)

type TxtHandler struct{}

func NewTxtHandler() *TxtHandler {
	return &TxtHandler{}
}

func (t *TxtHandler) Open(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
