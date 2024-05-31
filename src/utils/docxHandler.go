package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type DocxHandler struct{}

func NewDocxHandler() *DocxHandler {
	return &DocxHandler{}
}

func (d *DocxHandler) Open(path string) ([]byte, error) {
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

func (d *DocxHandler) ConvertToTxt(dir, file string) (string, error) {
	defer func() {
		os.Remove(filepath.Join(dir, file))
	}()
	newFile := strings.ReplaceAll(file, ".docx", "")
	cmd := exec.Command("docx2text", path.Join(dir, file), path.Join(dir, newFile))
	if err := cmd.Run(); err != nil {
		log.Printf("Error on convert pdf to txt: %v", err)
		return "", err
	}
	return strings.ReplaceAll(newFile, " ", "_"), nil
}
