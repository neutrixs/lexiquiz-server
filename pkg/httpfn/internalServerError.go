package httpfn

import (
	"fmt"
	"net/http"
)

func InternalServerError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte(fmt.Sprint(statusCode) + " " + http.StatusText(statusCode)))
	})
}