package workers

import (
	"context"
	"log"
	"strings"
	"time"

	"gitlab.com/bud.git/src/engine"
)

type WorkerListener struct {
	context.Context
	QuitChan chan bool
	*engine.Engine
	Workers  map[string]IWorker
	WorkerID string
}

func (w *WorkerListener) GetWorkerID() string {
	return w.WorkerID
}

func (w *WorkerListener) Spawn(ctx context.Context, id string, engine *engine.Engine, args ...any) IWorker {
	return &WorkerListener{
		Context:  ctx,
		WorkerID: id,
		QuitChan: make(chan bool),
		Engine:   engine,
		Workers:  w.Workers,
	}
}

func (w *WorkerListener) AddWorkers(workers map[string]IWorker) *WorkerListener {
	w.Workers = workers
	return w
}

func (w *WorkerListener) Run() {
	// startTime := time.Now()
	for {
		select {
		case <-w.Done():
			close(w.QuitChan)
			return
		case <-w.QuitChan:
			close(w.QuitChan)
			return
		default:
			time.Sleep(time.Millisecond * 100)
			// log.Printf("Worker running since %s", time.Since(startTime))
			continue
		}
	}
}

func (w *WorkerListener) Kill() error {
	w.QuitChan <- true
	return nil
}

func (w *WorkerListener) Call(args ...any) {
	w.AudioChan <- true
	question := <-w.AudioResponseChan
	cmd, err := w.ClassifySpeechCmd(question)
	if err != nil {
		log.Println("ERROR CLASSIFYING SPEECH COMMAND", err)
		return
	}
	cmd = strings.ToLower(cmd)
	cmd = strings.ReplaceAll(cmd, " ", "")
	cmd = strings.ReplaceAll(cmd, ":", "")

	for k, v := range w.Workers {
		// c := strings.Split(cmd, " ")
		log.Println("WORKER:", k)
		switch cmd {
		case k:
			log.Println(cmd)
			v.Call(question)
			return
		case "kill":
			if cmd == k {
				v.Kill()
				delete(w.Workers, k)
			}
			return
		default:
			log.Println("DEFAULTED THE COMMAND")
			if w.Workers["chat"] != nil {
				log.Println("CALLING CHAT WORKER")
				w.Workers["chat"].Call(question)
			} else {
				log.Println("NO WORKER FOUND")
				w.Speak("NO WORKER FOUND")
			}
			return
		}
	}
}
