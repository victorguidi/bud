package workers

import (
	"github.com/a-h/templ"
	wchat "gitlab.com/bud.git/src/workers/ChatWorker/view"
)

type ChatrHandler struct{}

func (a *WorkerChat) RegisterHandlers() {
	a.SERVEPAGE("/chatconfig", templ.Handler(wchat.WorkerChatView()))
}
