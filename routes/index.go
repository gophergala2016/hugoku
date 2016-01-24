package routes

import (
	"log"
	"net/http"
	"text/template"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"

	"github.com/gophergala2016/hugoku/store"
)

// Index is the Hugoku home page handler will redirect a non logged user to do the logging with Github
// or show a list of projectst and a form to add a project to a logged user,
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	session := sessions.GetSession(r)
	username := session.Get("username")

	if username == nil {
		t, err := template.ParseFiles("templates/index.html",
			"templates/partials/header.html",
			"templates/partials/footer.html",
		)
		if err != nil {
			log.Fatal("Error parsing the index page template: ", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal("Error executing the index page template: ", err)
		}
	} else {
		t, err := template.ParseFiles("templates/home.html",
			"templates/partials/header.html",
			"templates/partials/footer.html")
		if err != nil {
			log.Fatal("Error parsing the home page template: ", err)
		}
		log.Println("Recovering user data")
		user, err := store.GetUser(username.(string))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user)
		err = t.Execute(w, user)
		if err != nil {
			log.Fatal("Error executing the home page template: ", err)
		}
	}
}
