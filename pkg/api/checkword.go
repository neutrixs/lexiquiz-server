package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/words"
	"golang.org/x/exp/slices"
)

type checkWordResponse struct {
	Found bool `json:"found"`
}

type checkWordRequest struct {
	Word string `json:"word"`
}

func CheckWord(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
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
	data, err := words.WordsData.ReadFile("words.txt")
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
}