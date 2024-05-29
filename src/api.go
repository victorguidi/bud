package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/bud.git/src/engine"
	"gitlab.com/bud.git/src/workers"
)

type BudAPI struct {
	Workers map[string]workers.IWorker
	Engine  *engine.Engine
	Mux     *http.ServeMux
	Middleware
}

type DirBody struct {
	Dir []string `json:"dir"`
}

func (a *BudAPI) WithCors() {
	a.Middleware = chain(cors, defaultHandler)
}

func NewBudAPI(engine *engine.Engine) *BudAPI {
	api := &BudAPI{
		Mux:        http.NewServeMux(),
		Engine:     engine,
		Middleware: chain(defaultHandler),
	}
	return api
}

func (a *BudAPI) AddWorkers(workers map[string]workers.IWorker) *BudAPI {
	a.Workers = workers
	return a
}

func (a *BudAPI) Start(port string) {
	log.Printf("Http Server Listening on: %s", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, a.Mux))
}

func (a *BudAPI) RegisterHandlers() *BudAPI {
	// Dir ROUTES
	a.POST("/startrag", a.startragworker)
	a.POST("/stoprag", a.quitragworker)
	a.POST("/dir", a.dir)
	a.GET("/onedir/{dirname}", a.getOneDir)
	a.GET("/alldirs", a.getAllDirs)
	a.PUT("/dir/{dirname}", a.updateDir)
	a.DELETE("/dir", a.deleteDir)
	a.DELETE("/alldirs", a.deleteAllDirs)

	// Ask
	a.POST("/ask", a.dir)
	a.POST("/askbase", a.dir)
	a.POST("/askfile", a.dir)

	return a
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

func (a *BudAPI) startragworker(w http.ResponseWriter, r *http.Request) {
	go a.Workers["rag"].Run()
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Worker Dir started",
	})
}

func (a *BudAPI) quitragworker(w http.ResponseWriter, r *http.Request) {
	Workers["rag"].Stop()
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Stopping Work Dir",
	})
}

func (a *BudAPI) dir(w http.ResponseWriter, r *http.Request) {
	var dir DirBody
	json.NewDecoder(r.Body).Decode(&dir)

	cmd := []string{"new"}
	cmd = append(cmd, dir.Dir...)
	a.Workers["rag"].Call(cmd)

	json.NewEncoder(w).Encode(dir)
}

func (a *BudAPI) getOneDir(w http.ResponseWriter, r *http.Request) {
	dirname := r.PathValue("dirname")
	dir, err := a.Engine.SelectDir(dirname)
	if err != nil {
		http.Error(w, "Something Went Wrong", http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(dir)
}

func (a *BudAPI) getAllDirs(w http.ResponseWriter, r *http.Request) {
	dirs, err := a.Engine.SelectDirs()
	if err != nil {
		http.Error(w, "Something Went Wrong", http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(dirs)
}

func (a *BudAPI) updateDir(w http.ResponseWriter, r *http.Request) {}

func (a *BudAPI) deleteDir(w http.ResponseWriter, r *http.Request) {}

func (a *BudAPI) deleteAllDirs(w http.ResponseWriter, r *http.Request) {}

// mux.HandleFunc("/task/{id}/", func(w http.ResponseWriter, r *http.Request) {
//   id := r.PathValue("id")
//   fmt.Fprintf(w, "handling task with id=%v\n", id)
// })
