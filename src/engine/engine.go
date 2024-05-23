package engine

import (
	"log"
	"os"
	"path/filepath"
)

var dir = filepath.Join("testfiles")

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

// ProcessFiles read files (.pdf, .txt, .docx), generate Embeddings and send to a Postgres Vector Instance.
// Once Asked the Model will always have the files here as knowledge too.
func (e *Engine) ProcessFiles() {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Panic(err)
		}
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}
	for _, file := range files {
		log.Println(file)
	}
}

func (e *Engine) ProcessNews() {}
