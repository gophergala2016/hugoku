package routes

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"github.com/julienschmidt/httprouter"
	"github.com/google/go-github/github"
	"github.com/goincremental/negroni-sessions"

	"github.com/gophergala2016/hugoku/util/oauth"
	"github.com/gophergala2016/hugoku/store"
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
		log.Printf("client.Users.Get() faled with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var u store.User
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
