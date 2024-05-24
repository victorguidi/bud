package api

import (
	"fmt"
	"log"
	"net/http"
)

type BudAPI struct {
	Mux *http.ServeMux
}

func NewBudAPI() *BudAPI {
	api := &BudAPI{
		Mux: http.NewServeMux(),
	}
	return api
}

func (a *BudAPI) Start(port string) {
	log.Printf("Starting Server on Port: %s", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, a.Mux))
}

func (a *BudAPI) RegisterHandlers() {
	a.POST("/processfiles", a.processfiles)
}

func (a *BudAPI) GET(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("GET %s", path), handler)
}

func (a *BudAPI) POST(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("POST %s", path), handler)
}

func (a *BudAPI) PUT(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("PUT %s", path), handler)
}

func (a *BudAPI) DELETE(path string, handler http.HandlerFunc) {
	a.Mux.HandleFunc(fmt.Sprintf("DELETE %s", path), handler)
}

func (a *BudAPI) processfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "HELLOOOO")
}

// mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
//   id := r.PathValue("id")
//   fmt.Fprintf(w, "handling task with id=%v\n", id)
// })
