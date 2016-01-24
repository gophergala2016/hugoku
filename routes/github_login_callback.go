package routes

import (
	"log"
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"

	"github.com/gophergala2016/hugoku/store"
	"github.com/gophergala2016/hugoku/util/oauth"
)

//GithubLoginCallback Called by github after authorization is granted
func GithubLoginCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("GithubLoginCallback")
	state := r.FormValue("state")
	if state != oauth.RandomString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauth.RandomString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	// Convert the auth code into a token
	token, err := oauth.Conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get a client that uses the auth token
	oauthClient := oauth.Conf.Client(oauth2.NoContext, token)
	githubClient := github.NewClient(oauthClient)
	// Get the user info
	user, _, err := githubClient.Users.Get("")
	if err != nil {
		log.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	u, err := store.GetUser(*user.Login)
	if err != nil {
		log.Printf("store.GetUser failed with '%s'\n", err)
	}
	u.Username = *user.Login
	//	u.Email = *user.Email
	u.Token = *token
	u.AvatarURL = *user.AvatarURL
	err = store.SaveUser(u)

	session := sessions.GetSession(r)
	session.Set("username", *user.Login)

	log.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
