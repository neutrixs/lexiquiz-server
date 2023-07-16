package discordauth

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	database "github.com/neutrixs/lexiquiz-server/pkg/db"
)

func generateState(length int) string {
	const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890.-"
	var state string

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		randIndex := rand.Intn(len(possible))
		state += string([]rune(possible)[randIndex])
	}
	
	return state
}

func Login(w http.ResponseWriter, r *http.Request) {
	state := generateState(32)
	timestamp := time.Now().Unix()
	//TODO: implement checks if there was a duplicate state (literally impossible but, just in case)
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("prompt", "consent")
	query.Set("client_id", clientID)
	query.Set("scope", scopes)
	query.Set("redirect_uri", redirectURI)
	query.Set("state", state)
	qs := query.Encode()

	db := database.GetDB()
	_, err := db.Query("INSERT INTO discord_login (state, scopes, timestamp) VALUES (?, ?, ?)", state, scopes, timestamp)
	if err != nil {
		log.Println(err)
		statusCode := http.StatusInternalServerError
		w.WriteHeader(statusCode)
		w.Write([]byte(fmt.Sprint(statusCode) + " " + http.StatusText(statusCode)))
		return
	}
	
	http.Redirect(w, r, "https://discord.com/oauth2/authorize?" + qs, http.StatusFound)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}