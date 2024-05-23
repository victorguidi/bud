package engine

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/database"
	"gitlab.com/bud.git/src/utils"
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

	txtHandler := utils.NewTxtHandler()
	o := api.NewOllamaAPI()
	p := database.New()
	err = p.Initialize()
	if err != nil {
		log.Panic(err)
	}

	for _, file := range files {
		log.Println(file)

		fileBytes, err := txtHandler.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Panic(err)
		}

		e, err := o.GenerateEmbedding(context.Background(), string(fileBytes))
		if err != nil {
			log.Panic(err)
		}

		err = p.Save(file.Name(), string(fileBytes), e)
		if err != nil {
			log.Panic(err)
		}

	}

	emb, err := o.GenerateEmbedding(context.Background(), string([]byte("Victor")))
	if err != nil {
		log.Panic(err)
	}

	vectorTable, err := p.Retrieve(emb)
	if err != nil {
		log.Panic(err)
	}

	o.WithContext("How old is Victor", vectorTable.Text)
	call, err := o.SendMessageTo(context.Background())
	if err != nil {
		log.Panic(err)
	}

	log.Println(call)
}

func (e *Engine) ProcessNews() {}
