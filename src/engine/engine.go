package engine

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/database"
	"gitlab.com/bud.git/src/utils"
)

var (
	ollamaAPI = api.NewOllamaAPI()
	vectorDB  = database.New()
	dir       = filepath.Join("testfiles")
)

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

	err = e.EmbedFiles()
	if err != nil {
		panic(err)
	}

	emb, err := ollamaAPI.GenerateEmbedding(context.Background(), string([]byte("Victor")))
	if err != nil {
		log.Panic(err)
	}

	vectorTable, err := vectorDB.Retrieve(emb)
	if err != nil {
		log.Panic(err)
	}

	ollamaAPI.WithContext("What does the paper say about Bitcoin?", vectorTable.Text)
	call, err := ollamaAPI.SendMessageTo(context.Background())
	if err != nil {
		log.Panic(err)
	}

	log.Println(call)
}

func (e *Engine) EmbedFiles() error {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}

	var handler utils.IFiles
	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".pdf":
			pdf := utils.NewPDFHandler()
			pdf.ConvertToTxt(dir, file.Name())
			handler = pdf
		case ".docx":
			docx := utils.NewDocxHandler()
			docx.ConvertToTxt(dir, file.Name())
			handler = docx
		case ".txt":
			handler = utils.NewTxtHandler()
		default:
			return errors.New("FILE NOT RECOGNIZED")
		}

		fileBytes, err := handler.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Panic(err)
		}

		e, err := ollamaAPI.GenerateEmbedding(context.Background(), string(fileBytes))
		if err != nil {
			log.Panic(err)
		}

		err = vectorDB.Initialize()
		if err != nil {
			log.Panic(err)
		}

		err = vectorDB.Save(file.Name(), string(fileBytes), e)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func (e *Engine) ProcessNews() {}
