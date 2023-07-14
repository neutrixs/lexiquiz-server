package words

import (
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	switch strings.Split(r.URL.Path, "/")[0] {
	case "get":
		GetWord(w, r)
	case "check":
		CheckWord(w, r)
	default:
		http.NotFound(w, r)
	}
}