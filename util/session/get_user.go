package session

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/gophergala2016/hugoku/store"
)

// getUserFomSession gets the user from the usernane in the session on the http Request
func GetUser(r *http.Request) (store.User, error) {
	session := sessions.GetSession(r)
	username := session.Get("username")
	return store.GetUser(username.(string))
}
