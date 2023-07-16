package discordauth

import (
	"net/http"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/httpfn"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler

	switch strings.Split(r.URL.Path, "/")[0] {
	case "login":
		handler = httpfn.AllowMethod(http.MethodGet, http.HandlerFunc(Login))
	default:
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, r)
}