package workers

import (
	"context"
	"log"
	"strings"
	"time"

	workers "gitlab.com/bud.git/src/workers/RagWorker"

	"gitlab.com/bud.git/src/engine"
)

type WorkerListener struct {
	context.Context
	QuitChan chan bool
	*engine.Engine
	Workers  map[string]IWorker
	WorkerID string
	WState
}

// GetWorkerState implements IWorker.
func (w *WorkerListener) GetWorkerState() WState {
	return w.WState
}

func (w *WorkerListener) GetWorkerID() string {
	return w.WorkerID
}

func (w *WorkerListener) Spawn(ctx context.Context, id string, engine *engine.Engine, args ...any) IWorker {
	var r workers.RAG
	_ = r
	return &WorkerListener{
		Context:  ctx,
		WorkerID: id,
		QuitChan: make(chan bool),
		Engine:   engine,
		Workers:  w.Workers,
		WState:   ENABLED,
	}
}

func (w *WorkerListener) AddWorkers(workers map[string]IWorker) *WorkerListener {
	w.Workers = workers
	return w
}

func (w *WorkerListener) Run() {
	go w.LoadWhisper().Listen()
	w.WState = ENABLED
	log.Println("STARTING WORKER", w.WorkerID)
	// startTime := time.Now()
	for {
		select {
		case <-w.Done():
			close(w.QuitChan)
			return
		case <-w.QuitChan:
			log.Println("STOPPING WORKER ", w.WorkerID)
			w.StopListenerChan <- true
			return
		default:
			time.Sleep(time.Millisecond * 100)
			// log.Printf("Worker running since %s", time.Since(startTime))
			continue
		}
	}
}

func (w *WorkerListener) Stop() {
	w.WState = DISABLED
	w.QuitChan <- true
}

func (w *WorkerListener) Kill() error {
	log.Println("KILLING WORKER ", w.WorkerID)
	w.QuitChan <- true
	close(w.QuitChan)
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
	cmd = strings.ReplaceAll(cmd, "output: ", "")
	found := false

	for k, v := range w.Workers {
		switch strings.Contains(cmd, k) {
		case true:
			log.Println(cmd)
			v.Call([]string{question})
			found = true
			return
		case false:
			if cmd == "kill" {
				found = true
				log.Println("KILLING WORKER CALLED")
				// if cmd == k {
				//   v.Kill()
				//   delete(w.Workers, k)
				// }
			}
			continue
		}
	}
	if !found {
		log.Println("DEFAULTED THE COMMAND")
		if w.Workers["chat"] != nil {
			log.Println("CALLING CHAT WORKER")
			w.Workers["chat"].Call(question)
		} else {
			log.Println("NO WORKER FOUND")
			w.Speak("NO WORKER FOUND")
		}
	}
}
