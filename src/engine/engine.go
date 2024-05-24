package engine

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/database"
	"gitlab.com/bud.git/src/utils"
)

type Engine struct {
	EngineCommunicationPipe
	EngineFunctions
	EngineProperties
	ServerProperties
}

type EngineProperties struct {
	database.SqlDB
	database.IVectorDB
	api.OllamaAPI
}

type EngineFunctions struct {
	ProcessDirs func(paths []string, quit chan bool)
	ProcessFile func(path string)
	Crawler     func(sites []string)
}

type EngineCommunicationPipe struct {
	TriggerChan chan Trigger
	QuitChan    chan bool
}

type Trigger struct {
	Content  interface{}
	QuitChan chan bool
	Trigger  string
}

type DirTrigger struct {
	Dir string
}

type BaseTrigger struct {
	Question string
}

func New() *Engine {
	return &Engine{
		EngineCommunicationPipe{
			TriggerChan: make(chan Trigger),
			QuitChan:    make(chan bool),
		},
		EngineFunctions{},
		EngineProperties{
			OllamaAPI: *api.NewOllamaAPI(),
			IVectorDB: database.NewPostgresVectorDB(),
			SqlDB:     *database.NewSqlDB(),
		},
		ServerProperties{
			Host: "0.0.0.0",
			Port: "9876",
		},
	}
}

func (e *Engine) Run() {
	err := e.Initialize()
	if err != nil {
		log.Println("ERROR INITIALING THE VECTOR DB", err)
		return
	}

	err = e.Init()
	if err != nil {
		log.Println("ERROR INITIALING THE SQLITE DB", err)
		return
	}

	for {
		select {
		case cmd := <-e.TriggerChan:
			switch cmd.Trigger {

			case DIR.String():
				if content, ok := cmd.Content.(DirTrigger); ok {
					newDir := false
					if content.Dir != "" {
						_, err := os.Stat(content.Dir)
						if err != nil {
							if os.IsNotExist(err) {
								err := os.Mkdir(content.Dir, os.ModePerm)
								log.Println("CREATED DIR: ", content.Dir)
								if err != nil {
									log.Println("COULD NOT CREATE DIR:", content.Dir, err)
									return
								}
								newDir = true
							} else {
								log.Println("COULD NOT HANDLE DIR:", content.Dir, err)
								return
							}
						}
						err = e.InsertDirs(content.Dir)
						if err != nil {
							log.Println("COULD NOT SAVE THE DIR IN THE DB")
							return
						}
					}
					go e.ProcessDirs(cmd.QuitChan, newDir)
				}

			case ASKBASE.String():
				if question, ok := cmd.Content.(BaseTrigger); ok {
					go e.AskBase(question.Question)
				}

			default:
				continue
			}
		case <-e.QuitChan:
			fmt.Println("quit")
			return
		}
	}
}

func (e *Engine) Config() {}

func (e *Engine) ProcessDirs(quit chan bool, newDir bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second * 30)
			err := e.EmbedFiles(newDir)
			if err != nil {
				log.Println("ERROR CREATING THE EMBEDDINGS FOR THE FILES IN ONE OR MORE DIRECTORIES")
				return
			}
		}
	}
}

func (e *Engine) AskBase(question string) error {
	emb, err := e.GenerateEmbedding(context.Background(), string([]byte(question)))
	if err != nil {
		return err
	}

	vectorTable, err := e.Retrieve(emb.Embedding)
	if err != nil {
		return err
	}

	log.Printf("\n================\nVECTOR OUTPUT: %s\n================\n", vectorTable.Text)
	e.WithContext(question, vectorTable.Text)
	log.Println(question)
	call, err := e.SendMessageTo(context.Background())
	if err != nil {
		return err
	}

	audioEngine := NewAudioEngine()
	err = audioEngine.Speak(call.Response)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) EmbedFiles(newDir bool) error {
	dirs, err := e.SelectDirs()
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		log.Println("NO DIRS SAVED, PLEASE PROVIDE A DIR")
		return nil
	}

	for _, dir := range dirs {
		files, err := os.ReadDir(dir.Dir)
		if err != nil {
			log.Println("ERROR OPENING THE DIR:", dir, err)
			return err
		}

		var handler utils.IFiles
		var fileName string
		now := time.Now()

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				log.Println("ERROR GETTING FILE INFO:", file.Name(), err)
				continue
			}

			if !newDir {
				// Check if the file modification time is within the last 30 seconds
				if now.Sub(fileInfo.ModTime()) >= 60*time.Second {
					// Skip files older than 30 seconds
					continue
				}
			}

			switch filepath.Ext(file.Name()) {
			case ".pdf":
				pdf := utils.NewPDFHandler()
				newName, err := pdf.ConvertToTxt(dir.Dir, file.Name())
				if err != nil {
					return err
				}
				fileName = newName
				handler = pdf
			case ".docx":
				docx := utils.NewDocxHandler()
				newName, err := docx.ConvertToTxt(dir.Dir, file.Name())
				if err != nil {
					return err
				}
				fileName = newName
				handler = docx
			case ".txt":
				fileName = file.Name()
				handler = utils.NewTxtHandler()
			default:
				continue
			}

			fileBytes, err := handler.Open(filepath.Join(dir.Dir, fileName))
			if err != nil {
				return err
			}

			embedding, err := e.GenerateEmbedding(context.Background(), string(fileBytes))
			if err != nil {
				return err
			}

			err = e.Save(file.Name(), string(fileBytes), embedding.Embedding)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Engine) ProcessNews() {}
