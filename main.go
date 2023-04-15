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

	r := mux.NewRouter()
	r.HandleFunc("/api/getword", api.GetWord).Methods("GET")
	r.HandleFunc("/api/checkword", api.CheckWord).Methods("POST")
	
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(LISTEN, nil))
}