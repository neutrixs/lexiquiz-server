package discordauth

import (
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	switch strings.Split(r.URL.Path, "/")[0] {
	case "login":
		Login(w, r)
	default:
		http.NotFound(w, r)
	}
}