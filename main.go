package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// Startup instructions, Follow installation.md if you don't know what you are doing.:
// CREATE USER shitposting WITH LOGIN PASSWORD 'shitposting'
// create database shitposting
// create table items (item jsonb NOT NULL)

// Global variables.
var (
	db  *sql.DB
	pqd *pq.Driver
)

// APIHandler is the handler for /API. Holds the functions needed to handle cases of /API/{handle}
type APIHandler struct {
	Dirs map[string]*Directory
}

//Directory is to tree the handles.
type Directory struct {
	Function func(http.ResponseWriter, *http.Request)
	Dirs     map[string]*Directory
}

func (d *Directory) init() {
	// d.Function = func(http.ResponseWriter, *http.Request) {}
	d.Function = nil
	d.Dirs = make(map[string]*Directory)
}

// AddFunction is used to add functions to the Handler.
// Keep in mind that if you somehow add functions during runtime, they will not be kept for future runs.
func (h *APIHandler) AddFunction(dir []string, function func(http.ResponseWriter, *http.Request)) {
	var _, ok = h.Dirs[dir[0]]
	if !ok {
		h.Dirs[dir[0]] = &Directory{nil, make(map[string]*Directory)}
	}
	var Dir = h.Dirs[dir[0]]

	for _, v := range dir[1:] {
		_, ok = Dir.Dirs[v]
		if !ok {
			// Dir.Dirs[v] = Directory{func(http.ResponseWriter, *http.Request) {}, make(map[string]Directory)}
			d := Directory{}
			d.init()
			Dir.Dirs[v] = &Directory{nil, make(map[string]*Directory)}
		}
		Dir = Dir.Dirs[v]
	}
	Dir.Function = function

}

// ServeJSON serves back query results in json format.
func (h *APIHandler) ServeJSON(w http.ResponseWriter, r *http.Request) {
	var Query Item
	err := json.Unmarshal([]byte(r.PostFormValue("Query")), &Query)
	if err != nil {
		return
	}
	items := QueryLikeItem(Query)
	ToSend, err := json.Marshal(items)
	if err != nil {
		return
	}
	w.Write(ToSend)
}

// UpdateTag is used to update the tags of the same item within the db with the item received.
func (h *APIHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var Value Item
	err := json.Unmarshal([]byte(r.PostFormValue("item")), &Value)
	if err != nil {
		return
	}
	UpdateRow(Value)
	w.Write([]byte("OK"))
}

// ServeHTTP is the main point of entry to APIHandler.
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Args := strings.Split(r.RequestURI[5:], "/")
	dirs := h.Dirs[Args[0]]
	for _, v := range Args[1:] {
		if dirs.Dirs == nil {
			http.NotFound(w, r)
		}
		dirs = dirs.Dirs[v]
	}
	if dirs.Function == nil {
		http.NotFound(w, r)
	}
	Func := dirs.Function
	Func(w, r)
}

// ServeUser serves the main index page to the user.
func ServeUser(w http.ResponseWriter, r *http.Request) {
	p := "./www" + r.URL.Path
	if p == "./www/" {
		p = "./www/main.html"
	}
	http.ServeFile(w, r, p)

}

// ViewFile is used to serve the file.
func ViewFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	file := mux.Vars(r)["file"]
	Value := QuerySha(file)

	http.ServeFile(w, r, Value[0].File[0])

}

func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		return
	}

	r := mux.NewRouter()
	api := APIHandler{
		Dirs: make(map[string]*Directory),
	}
	api.AddFunction([]string{"JSON", "Query"}, api.ServeJSON)
	api.AddFunction([]string{"JSON", "UpdateTag"}, api.UpdateTag)
	r.PathPrefix("/API/").Handler(&api)
	r.PathPrefix("/view/{file}").HandlerFunc(ViewFile)
	r.PathPrefix("/").HandlerFunc(ServeUser)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go initDB()

	srv.ListenAndServe()

}
