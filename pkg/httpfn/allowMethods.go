package httpfn

import (
	"fmt"
	"net/http"
)

func AllowMethod(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			status := http.StatusMethodNotAllowed
			w.WriteHeader(status)
			w.Write([]byte(fmt.Sprint(status) + " " + http.StatusText(status)))
			return
		}

		h.ServeHTTP(w, r)
	})
}