package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/neutrixs/lexiquiz-server/pkg/api"
	"github.com/neutrixs/lexiquiz-server/pkg/env"
)

func disableDirectoryListing(pathname string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target, err := os.Stat(path.Join(pathname, r.URL.Path))
		if err == nil && target.IsDir() {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

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

	http.Handle("/static/", http.StripPrefix("/static/", disableDirectoryListing(STATIC, http.FileServer(http.Dir(STATIC)))))
	http.Handle("/api/", http.StripPrefix("/api/", http.HandlerFunc(api.Handle)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, INDEX)
	})
	log.Fatal(http.ListenAndServe(LISTEN, nil))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}