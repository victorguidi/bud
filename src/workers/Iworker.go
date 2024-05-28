package workers

import (
	"context"

	"gitlab.com/bud.git/src/engine"
)

type IWorker interface {
	GetWorkerID() string
	Spawn(context.Context, string, *engine.Engine, ...any) IWorker
	Kill() error
	Call(...any)
	Run()
}
