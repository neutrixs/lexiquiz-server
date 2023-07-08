package api

import (
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/words"
	"golang.org/x/exp/slices"
)

func GetWord(w http.ResponseWriter, r *http.Request) {
	data, err := words.WordsData.ReadFile("common.txt")
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
}