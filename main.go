package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neutrixs/lexiquiz-server/pkg/api"
	"github.com/neutrixs/lexiquiz-server/pkg/env"
)

func main() {
	LISTEN, err := env.Get("LISTEN")
	if err != nil {
		LISTEN = ":8080"
	}

	STATIC, err := env.Get("STATIC_DIR")
	if err != nil {
		log.Fatal(err)
	}

	INDEX, err := env.Get("INDEX_HTML")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/getword", api.GetWord).Methods("GET")
	r.HandleFunc("/api/checkword", api.CheckWord).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir(STATIC))))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, INDEX)
	})
	
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(LISTEN, nil))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}