package routes

import (
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

// About shows info about the project, team  etc ...
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/about.html",
		"templates/partials/header.html",
		"templates/partials/footer.html",
	)
	if err != nil {
		log.Fatal("Error parsing the about page template")
	}
	t.Execute(w, nil)
}
