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
	WState
}

type trigger struct {
	Question string
}

func (w *WorkerChat) GetWorkerID() string {
	return w.WorkerID
}

func (w *WorkerChat) GetWorkerState() WState {
	return w.WState
}

func (w *WorkerChat) Spawn(ctx context.Context, id string, engine *engine.Engine, args ...any) IWorker {
	return &WorkerChat{
		Context:     ctx,
		WorkerID:    id,
		TriggerChan: make(chan trigger),
		QuitChan:    make(chan bool),
		Engine:      engine,
		WState:      ENABLED,
	}
}

func (w *WorkerChat) Run() {
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

func (w *WorkerChat) Stop() {
	w.WState = DISABLED
	w.QuitChan <- true
}

func (w *WorkerChat) Kill() error {
	log.Println("KILLING WORKER ", w.WorkerID)
	close(w.TriggerChan)
	w.QuitChan <- true
	close(w.QuitChan)
	return nil
}

func (w *WorkerChat) Call(args ...any) {
	if w.String() == "on" {
		for _, a := range args {
			if q, ok := a.([]string); ok {
				w.TriggerChan <- trigger{
					Question: parseQuestion(q),
				}
			} else if a, ok := a.(string); ok {
				w.TriggerChan <- trigger{
					Question: a,
				}
			}
		}
	} else {
		log.Printf("WORKER %s IS NOT ENABLED\n", w.WorkerID)
	}
}

func (w *WorkerChat) askLLM(question string) error {
	log.Println(question)
	w.WithTokens(100) // This should be dynamic and have a way of changing with CLI and API
	w.PromptFormater(api.DEFAULTPROMPT, map[string]string{
		"Input": question,
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
		question.WriteString(fmt.Sprintf("%s ", p))
	}
	return question.String()
}
