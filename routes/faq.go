package routes

import (
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

// FAQ shows frequently asqued questions about the project, team  etc ...
func FAQ(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/faq.html",
		"templates/partials/header.html",
		"templates/partials/footer.html")
	if err != nil {
		log.Fatal("Error parsing the FAQ page template: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("Error executing the FAQ page template: ", err)
	}
}
