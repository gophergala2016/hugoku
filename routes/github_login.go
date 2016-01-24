package routes

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"

	"github.com/gophergala2016/hugoku/util/oauth"
)

// GithubLogin redirects the user to github to handle the Oauth2 authentication
func GithubLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("GithubLogin")

	url := oauth.Conf.AuthCodeURL(oauth.RandomString, oauth2.AccessTypeOnline)
	log.Println("Redirecting the user to github login")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
