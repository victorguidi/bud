package engine

import (
	"net/http"

	"github.com/a-h/templ"
)

func (e *Engine) Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}
