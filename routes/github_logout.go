package routes

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
)

func GithubLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessions.GetSession(r)
	session.Delete("username")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
