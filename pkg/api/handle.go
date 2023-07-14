package api

import (
	"net/http"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/api/words"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler

	switch strings.Split(r.URL.Path, "/")[0] {
	case "words" :
		handler = http.StripPrefix("words/", http.HandlerFunc(words.Handle))
	default:
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}