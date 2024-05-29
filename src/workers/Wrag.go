package workers

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/database"
	"gitlab.com/bud.git/src/engine"
	"gitlab.com/bud.git/src/utils"
)

type WorkerRag struct {
	context.Context
	WorkerID    string
	TriggerChan chan ragtrigger
	QuitChan    chan bool
	*engine.Engine
	WState
	database.ISqlDB[WorkerRagConfig]
}

type WorkerRagConfig struct{}

type ragtrigger struct {
	Question string
}

func (w *WorkerRag) GetWorkerID() string {
	return w.WorkerID
}

func (w *WorkerRag) GetWorkerState() WState {
	return w.WState
}

func (w *WorkerRag) Spawn(ctx context.Context, id string, engine *engine.Engine, args ...any) IWorker {
	return &WorkerRag{
		Context:     ctx,
		WorkerID:    id,
		TriggerChan: make(chan ragtrigger),
		QuitChan:    make(chan bool),
		Engine:      engine,
		WState:      ENABLED,
		ISqlDB:      database.NewSqlDB[WorkerChatConfig](),
	}
}

func (w *WorkerRag) Run() {
	w.WState = ENABLED
	log.Println("STARTING WORKER", w.WorkerID)
	// startTime := time.Now()
	for {
		select {
		case <-w.Done():
			close(w.TriggerChan)
			close(w.QuitChan)
			return
		case <-w.QuitChan:
			log.Println("STOPPING WORKER ", w.WorkerID)
			return
		case t := <-w.TriggerChan:
			err := w.AskBase(t.Question)
			if err != nil {
				log.Println("ERROR ASKING THE LLM", err)
			}
		default:
			time.Sleep(time.Second * 5)
			err := w.EmbedFiles(false)
			if err != nil {
				log.Println("ERROR CREATING THE EMBEDDINGS FOR THE FILES IN ONE OR MORE DIRECTORIES")
				return
			}
			continue
		}
	}
}

func (w *WorkerRag) Stop() {
	w.WState = DISABLED
	w.QuitChan <- true
}

func (w *WorkerRag) Kill() error {
	log.Println("KILLING WORKER ", w.WorkerID)
	close(w.TriggerChan)
	w.QuitChan <- true
	close(w.QuitChan)
	return nil
}

func (w *WorkerRag) Call(args ...any) {
	if w.String() == "on" {
		for _, a := range args {
			if cmd, ok := a.([]string); ok {
				switch cmd[0] {
				case "new":
					for _, dir := range cmd[1:] {
						err := w.Insert(dir)
						if err != nil {
							log.Println("ERROR INSERTING DIR: ", dir)
						}
					}
					w.EmbedFiles(true)
				case "chat":
					w.TriggerChan <- ragtrigger{
						Question: strings.Join(cmd[1:], " "),
					}
				default:
					log.Println(cmd)
					w.TriggerChan <- ragtrigger{
						Question: strings.Join(cmd, " "),
					}
				}
			}
		}
	} else {
		log.Printf("WORKER %s IS NOT ENABLED\n", w.WorkerID)
	}
}

// EmbedFiles will look at dirs that are stored at the SQLite database, if the dir is new it will embed all files inside, else it will embed only new files
// New files are considered if they were created or edit in the last 60 seconds
func (w *WorkerRag) EmbedFiles(newDir bool) error {
	dirs, err := w.SelectDirs()
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		log.Println("NO DIRS SAVED, PLEASE INSERT A NEW DIR")
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
				log.Println("WORKER ONLY FOUND DIRS FROM THE BASE, DEFAULT TO NEW FILES ONLY")
				if now.Sub(fileInfo.ModTime()) >= 60*time.Second {
					// Skip files older than 30 seconds
					log.Println("OLD FILE", file.Name())
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

			embedding, err := w.GenerateEmbedding(w.Context, string(fileBytes))
			if err != nil {
				return err
			}

			err = w.Save(file.Name(), string(fileBytes), embedding.Embedding)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SECTION: Talk to the model
// AskBase will prompt the model but it will use the files that it could find on the directores provided in the DIR section.
func (w *WorkerRag) AskBase(question string) error {
	emb, err := w.GenerateEmbedding(w.Context, string([]byte(question)))
	if err != nil {
		return err
	}

	vectorTable, err := w.Retrieve(emb.Embedding)
	if err != nil {
		return err
	}

	log.Printf("\n================\nVECTOR OUTPUT: %s\n================\n", vectorTable.Text)
	w.PromptFormater(api.DEFAULTRAGPROMPT, map[string]string{
		"Context": vectorTable.Text,
		"Input":   question,
	})

	w.WithTokens(150)
	call, err := w.SendMessageTo(w.Context)
	if err != nil {
		return err
	}

	audioEngine := engine.NewAudioEngine()
	err = audioEngine.Speak(call.Response)
	if err != nil {
		return err
	}
	return nil
}
