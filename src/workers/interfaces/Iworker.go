package workers

import (
	"context"

	"gitlab.com/bud.git/src/engine"
)

type WState string

// Constants for common HTTP methods
const (
	ENABLED  WState = "on"
	DISABLED WState = "off"
)

func (m WState) String() string {
	return string(m)
}

type IWorker interface {
	GetWorkerID() string
	GetWorkerState() WState
	Spawn(context.Context, string, *engine.Engine, ...any) IWorker
	Run()
	Stop()
	Kill() error
	Call(...any)
}
