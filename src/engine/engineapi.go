package engine

import (
	"fmt"
	"log"
	"net/http"
)

type IAPI interface {
	RegisterHandlers()
}

type BudAPI struct {
	Mux *http.ServeMux
	Middleware
}

func (a *BudAPI) ExtendRoutes(apis ...IAPI) *BudAPI {
	for _, api := range apis {
		api.RegisterHandlers()
	}
	return a
}

func (a *BudAPI) WithCors() {
	a.Middleware = chain(cors, defaultHandler)
}

func NewBudAPI() *BudAPI {
	return &BudAPI{
		Mux:        http.NewServeMux(),
		Middleware: chain(defaultHandler),
	}
}

func (a *BudAPI) Start(port string) {
	log.Printf("Http Server Listening on: %s", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, a.Mux))
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func chain(middleware ...Middleware) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		for _, m := range middleware {
			handler = m(handler)
		}
		return handler
	}
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func defaultHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func (a *BudAPI) GET(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("GET %s", path), a.Middleware(handler))
}

func (a *BudAPI) POST(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("POST %s", path), a.Middleware(handler))
}

func (a *BudAPI) PUT(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("PUT %s", path), a.Middleware(handler))
}

func (a *BudAPI) DELETE(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("DELETE %s", path), a.Middleware(handler))
}
