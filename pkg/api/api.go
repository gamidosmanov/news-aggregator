package api

import (
	"encoding/json"
	"net/http"
	"news-aggregator/pkg/storage"
	"strconv"

	"github.com/gorilla/mux"
)

// API encapsulates api internals
type API struct {
	db     *storage.DB
	router *mux.Router
}

// New returns new API object provided DB object
func New(db *storage.DB) *API {
	api := API{
		db:     db,
		router: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// Router get api router
func (api *API) Router() *mux.Router {
	return api.router
}

// endpoints registers handlers
func (api *API) endpoints() {
	api.router.HandleFunc("/news/{n}", api.getPosts).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// getPosts fetches n posts from database
func (api *API) getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	posts, err := api.db.Posts(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
