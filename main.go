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

// CREATE USER shitposting WITH LOGIN PASSWORD 'hunter2'
// create database shitposting
// create table items (item jsonb NOT NULL)

// select * from items where item->'tags' @> '["tag"]';

// insert into items (item) values ('{"file": "/location", "tags": ["1"], "thumbnail": "thumbnail/location", "sha1":"9977"}');

var db *sql.DB
var pqd *pq.Driver

// Item is the main json used to read from and write to the database.
type Item struct {
	File      []string `json:"location"`
	Thumbnail string   `json:"thumbnail"`
	Tags      []string `json:"tags"`
	Sha1      string   `json:"sha1"`
}

// APIHandler is the handler for /API. Holds the functions needed to handle cases of /API/{handle}
type APIHandler struct {
	Functions map[string]func(http.ResponseWriter, *http.Request)
}

// AddFunction is experimental and might be used for adding new handlers during runtime
// Keep in mind that functions added during runtime will not be kept for future runs.
func (h *APIHandler) AddFunction(name string, function func(http.ResponseWriter, *http.Request)) {
	h.Functions[name] = function
}

// ServeJSON serves back query results in json format.
func (h *APIHandler) ServeJSON(w http.ResponseWriter, r *http.Request) {
	var Query Item
	err := json.Unmarshal([]byte(r.PostFormValue("Query")), &Query)
	if err != nil {
		return
	}
	rows, err := QueryLikeItem(Query)
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var received Item
		var testing string
		err := rows.Scan(&testing)
		if err != nil {
			// log.Println(err)
		}
		json.Unmarshal([]byte(testing), &received)
		items = append(items, received)
	}
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

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Args := mux.Vars(r)
	Func, ok := h.Functions[Args["handle"]+Args["method"]]
	if ok {
		Func(w, r)
	}
}

// ServeUser serves the main index page to the user.
func ServeUser(w http.ResponseWriter, r *http.Request) {
	var Folder string
	var file bool
	Folder = string([]rune(r.URL.Path)[:strings.LastIndex(r.URL.Path, "/")])
	for _, x := range config.Folders {
		if Folder == x {
			file = true
		}
	}
	if file {
		http.ServeFile(w, r, r.URL.Path)

	} else {
		p := "." + r.URL.Path
		if p == "./" {
			p = "./main.html"
		}
		http.ServeFile(w, r, p)
	}
}

// func ServeThumbnail(w http.ResponseWriter, r * http.Request){

// }

func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		return
	}

	r := mux.NewRouter()
	api := APIHandler{
		Functions: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	api.AddFunction("JSONQuery", api.ServeJSON)
	api.AddFunction("JSONUpdateTag", api.UpdateTag)
	r.PathPrefix("/API/{handle}/{method}").Handler(&api)
	r.PathPrefix("/").HandlerFunc(ServeUser)
	// r.PathPrefix("/thumbnail/").HandlerFunc(ServeThumbnail)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go initDB()

	srv.ListenAndServe()

}
