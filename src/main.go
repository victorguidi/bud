package main

// TODO: Create a simple Crawler For Websites
// TODO: Simple Frontend Client that is spawned with service (Talk | search docs)
// TODO: Improve memory usage

import (
	"context"

	"gitlab.com/bud.git/src/engine"
	"gitlab.com/bud.git/src/workers"
)

var Workers = make(map[string]workers.IWorker)

func main() {
	ctx := context.Background()
	bud := engine.New()

	bud.Run()                                                // Start the Engine
	go bud.Listen()                                          // Start the Augio Engine
	go NewServerEngine(ctx, "0.0.0.0", "9876").StartServer() // Start the Engine Socket

	// Register Workers
	registerWorkes(
		new(workers.WorkerChat).Spawn(ctx, "chat", bud),
		new(workers.WorkerListener).AddWorkers(Workers).Spawn(ctx, "listen", bud),
	)

	NewBudAPI(bud).RegisterHandlers().Start("9875")
}

func registerWorkes(workers ...workers.IWorker) {
	for _, w := range workers {
		Workers[w.GetWorkerID()] = w
		go w.Run()
	}
}
