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

type PDFHandler struct{}

func NewPDFHandler() *PDFHandler {
	return &PDFHandler{}
}

func (p *PDFHandler) Open(path string) ([]byte, error) {
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

func (p *PDFHandler) ConvertToTxt(dir, file string) (string, error) {
	defer func() {
		os.Remove(filepath.Join(dir, file))
	}()
	newFile := strings.ReplaceAll(file, ".pdf", ".txt")
	cmd := exec.Command("pdftotext", path.Join(dir, file), path.Join(dir, newFile))
	if err := cmd.Run(); err != nil {
		log.Printf("Error on convert pdf to txt: %v", err)
		return "", err
	}
	return strings.ReplaceAll(newFile, " ", "_"), nil
}
