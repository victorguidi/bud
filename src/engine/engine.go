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

type Engine struct {
	Question string
}

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

	// emb, err := ollamaAPI.GenerateEmbedding(context.Background(), string([]byte("Where is the Capital of New Zealand")))
	emb, err := ollamaAPI.GenerateEmbedding(context.Background(), string([]byte(e.Question)))
	if err != nil {
		log.Panic(err)
	}

	vectorTable, err := vectorDB.Retrieve(emb.Embedding)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("\n================\nVECTOR OUTPUT: %s\n================\n", vectorTable.Text)
	ollamaAPI.WithContext(e.Question, vectorTable.Text)
	call, err := ollamaAPI.SendMessageTo(context.Background())
	if err != nil {
		log.Panic(err)
	}

	audioEngine := NewAudioEngine()
	err = audioEngine.Speak(call.Response)
	if err != nil {
		log.Panic(err)
	}
}

func (e *Engine) EmbedFiles() error {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Panic(err)
	}

	var handler utils.IFiles
	var fileName string
	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".pdf":
			pdf := utils.NewPDFHandler()
			newName, err := pdf.ConvertToTxt(dir, file.Name())
			if err != nil {
				return err
			}
			fileName = newName
			handler = pdf
		case ".docx":
			docx := utils.NewDocxHandler()
			newName, err := docx.ConvertToTxt(dir, file.Name())
			if err != nil {
				return err
			}
			fileName = newName
			handler = docx
		case ".txt":
			fileName = file.Name()
			handler = utils.NewTxtHandler()
		default:
			return errors.New("FILE NOT RECOGNIZED")
		}

		fileBytes, err := handler.Open(filepath.Join(dir, fileName))
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

		err = vectorDB.Save(file.Name(), string(fileBytes), e.Embedding)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func (e *Engine) ProcessNews() {}
