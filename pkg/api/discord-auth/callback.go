package discordauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	database "github.com/neutrixs/lexiquiz-server/pkg/db"
	"github.com/neutrixs/lexiquiz-server/pkg/httpfn"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	query := r.URL.Query()

	var (
		callbackError = query.Get("error")
		code = query.Get("code")
		state = query.Get("state")
	)
	if callbackError != "" {
		w.Write([]byte("<script>window.close();</script>"))
		return
	}
	if code == "" || state == "" {
		httpfn.BadRequest().ServeHTTP(w, r)
		return
	}

	rows, err := db.Query("SELECT timestamp, scopes FROM discord_login WHERE state = ?", state)
	if err != nil {
		fmt.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
		return
	}

	var (
		exists bool
		timestamp int
		scopes string
		currentTime = time.Now().Unix()
	)

	for rows.Next() {
		exists = true
		err := rows.Scan(&timestamp, &scopes)
		if err != nil {
			fmt.Println(err)
			httpfn.InternalServerError().ServeHTTP(w, r)
			return
		}
	}

	rows.Close()

	if !exists {
		statusCode := http.StatusNotFound
		w.WriteHeader(statusCode)
		w.Write([]byte("state not found"))
		return
	}

	if currentTime - int64(timestamp) > 300 {
		statusCode := http.StatusForbidden
		w.WriteHeader(statusCode)
		w.Write([]byte("the authentication has expired"))
		return
	}

	param := url.Values{}
	param.Set("grant_type", "authorization_code")
	param.Set("code", code)
	param.Set("redirect_uri", redirectURI)
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	request, err := http.NewRequest("POST", APIEndpoint + "/oauth2/token", strings.NewReader(param.Encode()))
	if err != nil {
		log.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic " + auth)

	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
		return
	}

	token := AccessTokenResponse{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
	}

	if !compareScopes(scopes, token.Scope) {
		statusCode := http.StatusForbidden
		w.WriteHeader(statusCode)
		w.Write([]byte(fmt.Sprint(statusCode) + " " + http.StatusText(statusCode)))
		return
	}

	var (
		cookie http.Cookie
		loginProvider = "discord"
		currentTimestamp = time.Now().Unix()
	)

	cookie.Name = "state"
	cookie.Value = state
	cookie.Expires = time.Now().Add(36500 * 24 * time.Hour)

	_, err = db.Query(
		"INSERT INTO user_auth (state, login_provider, refresh_token, timestamp) VALUES (?, ?, ?, ?)",
		state, loginProvider, token.RefreshToken, currentTimestamp,
	)
	if err != nil {
		log.Println(err)
		httpfn.InternalServerError().ServeHTTP(w, r)
		return
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<script>window.close();</script>"))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}