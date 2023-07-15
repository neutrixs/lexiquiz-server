package api

import (
	"net/http"
	"strings"

	discordauth "github.com/neutrixs/lexiquiz-server/pkg/api/discord-auth"
	"github.com/neutrixs/lexiquiz-server/pkg/api/words"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler

	switch strings.Split(r.URL.Path, "/")[0] {
	case "words" :
		handler = http.StripPrefix("words/", http.HandlerFunc(words.Handle))
	case "discord-auth":
		handler = http.StripPrefix("discord-auth/", http.HandlerFunc(discordauth.Handle))
	default:
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}