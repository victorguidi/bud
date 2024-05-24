package api

import (
	"fmt"
	"log"
	"net/http"
)

type BudAPI struct {
	Mux *http.ServeMux
	Middleware
}

// type Handler struct{}
//
// func (h *Handler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	h.ServerHTTP(w, r)
// }

func NewBudAPI() *BudAPI {
	api := &BudAPI{
		Mux:        http.NewServeMux(),
		Middleware: chain(defaultHandler),
	}
	return api
}

func (a *BudAPI) WithCors() {
	a.Middleware = chain(cors, defaultHandler)
}

func (a *BudAPI) Start(port string) {
	log.Printf("Starting Server on Port: %s", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, a.Mux))
}

func (a *BudAPI) RegisterHandlers() {
	// Dir ROUTES
	a.POST("/dir", a.processfiles)
	a.GET("/onedir/{dirname}", a.processfiles)
	a.GET("/alldirs", a.processfiles)
	a.PUT("/dir/{dirname}", a.processfiles)
	a.DELETE("/dir", a.processfiles)
	a.DELETE("/alldirs", a.processfiles)

	// Ask
	a.POST("/ask", a.processfiles)
	a.POST("/askbase", a.processfiles)
	a.POST("/askfile", a.processfiles)
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
	a.Mux.HandleFunc(fmt.Sprintf("GET %s", path), defaultHandler(handler))
}

func (a *BudAPI) POST(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("POST %s", path), defaultHandler(handler))
}

func (a *BudAPI) PUT(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("PUT %s", path), defaultHandler(handler))
}

func (a *BudAPI) DELETE(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("DELETE %s", path), defaultHandler(handler))
}

func (a *BudAPI) processfiles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HELLOOOO")
}

// mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
//   id := r.PathValue("id")
//   fmt.Fprintf(w, "handling task with id=%v\n", id)
// })
