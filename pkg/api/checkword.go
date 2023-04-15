package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/env"
	"golang.org/x/exp/slices"
)

type checkWordResponse struct {
	Found bool `json:"found"`
}

type checkWordRequest struct {
	Word string `json:"word"`
}

func CheckWord(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status := http.StatusBadRequest
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
		return
	}

	var body checkWordRequest
	if err := json.Unmarshal(requestBody, &body); err != nil {
		status := http.StatusBadRequest
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
		return
	}

	word := body.Word
	path, _ := env.Get("WORDS_PATH")
	data, err := os.ReadFile(path)
	if err != nil {
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
		return
	}

	var response checkWordResponse
	words := strings.Split(string(data), "\n")
	response.Found = slices.Contains(words, word)

	encodedResponse, err := json.Marshal(response)
	if err != nil {
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
		return
	}

	w.Write(encodedResponse)
	w.Header().Add("Content-Type", "application/json")
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	_, err := env.Get("WORDS_PATH")
	if err != nil {
		log.Fatal(err)
	}
}