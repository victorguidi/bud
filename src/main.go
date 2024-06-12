package main

// TODO: Create a simple Crawler For Websites
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
// TODO: Improve memory usage
// FIX: Memory allocation for Whisper is quite weird

import (
	"context"

	"gitlab.com/bud.git/src/engine"
	"gitlab.com/bud.git/src/workers"
	wchat "gitlab.com/bud.git/src/workers/ChatWorker"
)

var Workers = make(map[string]workers.IWorker)

func main() {
	ctx := context.Background()
	bud := engine.New()

	bud.Run() // Start the Engine

	go NewServerEngine(ctx, "0.0.0.0", "9876").StartServer() // Start the Engine Socket

	// Register Workers
	go registerWorkes(
		new(wchat.WorkerChat).Spawn(ctx, "chat", bud),
		// new(workers.WorkerRag).Spawn(ctx, "rag", bud),
		new(workers.WorkerListener).AddWorkers(Workers).Spawn(ctx, "listen", bud),
	)

	// Start the HTTP SERVER
	bud.WithCors().Start("9875")
}

func registerWorkes(workers ...workers.IWorker) {
	for _, w := range workers {
		Workers[w.GetWorkerID()] = w
		go w.Run()
	}
}
