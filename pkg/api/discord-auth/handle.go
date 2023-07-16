package discordauth

import (
	"net/http"
	"strings"

	"github.com/neutrixs/lexiquiz-server/pkg/env"
	"github.com/neutrixs/lexiquiz-server/pkg/httpfn"
	"golang.org/x/exp/slices"
)

var (
	scopes string
	redirectURI string
	clientID string
	clientSecret string
	APIEndpoint = "https://discord.com/api/v10"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
}

func compareScopes(scopes1 string, scopes2 string) bool {
	s1 := strings.Split(scopes1, " ")
	s2 := strings.Split(scopes2, " ")

	slices.Sort(s1)
	slices.Sort(s2)

	return strings.Join(s1, " ") == strings.Join(s2, " ")
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler

	switch strings.Split(r.URL.Path, "/")[0] {
	case "login":
		handler = httpfn.AllowMethod(http.MethodGet, http.HandlerFunc(Login))
	case "callback":
		handler = httpfn.AllowMethod(http.MethodGet, http.HandlerFunc(Callback))
	default:
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, r)
}

func init() {
	scopes, _ = env.Get("DISCORD_SCOPES")
	redirectURI, _ = env.Get("DISCORD_REDIRECT_URI")
	clientID, _ = env.Get("DISCORD_CLIENT_ID")
	clientSecret, _ = env.Get("DISCORD_CLIENT_SECRET")
}