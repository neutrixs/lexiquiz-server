package api

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/env"
	"golang.org/x/exp/slices"
)

func GetWord(w http.ResponseWriter, r *http.Request) {
	path, _ := env.Get("COMMON_WORDS_PATH")
	data, err := os.ReadFile(path)
	if err != nil {
		status := http.StatusInternalServerError
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
		return
	}

	excludeQuery := r.URL.Query().Get("exclude")
	exclude := strings.Split(excludeQuery, ",")
	words := strings.Split(string(data), "\n")
	var selectedWord string

	for selectedWord == "" {
		randomIndex := rand.Intn(len(words))
		word := words[randomIndex]

		if !slices.Contains(exclude, word) {
			selectedWord = word
		}
	}

	w.Write([]byte(selectedWord))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	_, err := env.Get("COMMON_WORDS_PATH")
	if err != nil {
		log.Fatal(err)
	}
}