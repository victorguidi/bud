package workers

import (
	"net/http"

	"gitlab.com/bud.git/src/workers/ChatWorker/wchat"
)

func (a *WorkerChat) RegisterHandlers() {
	a.GET("/chatconfig", a.handleShowChatConfig)
}

func (a *WorkerChat) handleShowChatConfig(w http.ResponseWriter, r *http.Request) {
	err := a.Render(w, r, wchat.WorkerChatView())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
