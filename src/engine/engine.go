package engine

import (
	"context"
	"log"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/database"
)

type Engine struct {
	EngineProperties[Engine]
	context.Context
	AudioEngine
	BudAPI
}

type EngineProperties[T any] struct {
	database.ISqlDB[T]
	database.IVectorDB
	api.OllamaAPI
}

func New() *Engine {
	return &Engine{
		EngineProperties[Engine]{
			OllamaAPI: *api.NewOllamaAPI(),
			IVectorDB: database.NewPostgresVectorDB(),
			ISqlDB:    database.NewSqlDB[Engine](),
		},
		context.Background(),
		AudioEngine{
			AudioChan:         make(chan bool),
			AudioResponseChan: make(chan string),
			StopListenerChan:  make(chan bool),
		},
		*NewBudAPI(),
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
}

func (e *Engine) Config() {}

func (e *Engine) ClassifySpeechCmd(cmd string) (string, error) {
	e.PromptFormater(api.DEFAULTCLASSIFIER, map[string]string{
		"Input": cmd,
	})

	e.WithTokens(100) // Reduce the ammount of tokens to 100 only
	call, err := e.SendMessageTo(context.Background())
	if err != nil {
		return "", err
	}

	log.Println("CLASSIFIED RESPONSE: ", call.Response)
	return call.Response, nil
}

func (e *Engine) ProcessNews() {}
