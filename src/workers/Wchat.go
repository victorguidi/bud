package workers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.com/bud.git/src/api"
	"gitlab.com/bud.git/src/engine"
)

type WorkerChat struct {
	context.Context
	WorkerID    string
	TriggerChan chan trigger
	QuitChan    chan bool
	*engine.Engine
}

type trigger struct {
	Question string
}

func (w *WorkerChat) GetWorkerID() string {
	return w.WorkerID
}

func (w *WorkerChat) Spawn(ctx context.Context, id string, engine *engine.Engine, args ...any) IWorker {
	return &WorkerChat{
		Context:     ctx,
		WorkerID:    id,
		TriggerChan: make(chan trigger),
		QuitChan:    make(chan bool),
		Engine:      engine,
	}
}

func (w *WorkerChat) Run() {
	// startTime := time.Now()
	for {
		select {
		case <-w.Done():
			return
		case <-w.QuitChan:
			return
		case t := <-w.TriggerChan:
			err := w.askLLM(t.Question)
			if err != nil {
				log.Println("ERROR ASKING THE LLM", err)
			}
		default:
			time.Sleep(time.Millisecond * 100)
			// log.Printf("Worker running since %s", time.Since(startTime))
			continue
		}
	}
}

func (w *WorkerChat) Kill() error {
	close(w.TriggerChan)
	close(w.QuitChan)
	return nil
}

func (w *WorkerChat) Call(args ...any) {
	for _, a := range args {
		if a, ok := a.([]string); ok {
			w.TriggerChan <- trigger{
				Question: parseQuestion(a),
			}
		}
	}
}

func (w *WorkerChat) askLLM(question string) error {
	w.PromptFormater(api.DEFAULTPROMPT, map[string]string{
		"question": question,
	})

	call, err := w.SendMessageTo(context.Background())
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

func parseQuestion(prompt []string) string {
	var question strings.Builder
	for _, p := range prompt {
		question.WriteString(fmt.Sprintf("%s", p))
	}
	return question.String()
}
