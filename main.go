package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("Hello hugoku!")
	Serve()
}

// Serve set the route handlers and serve
func Serve() {

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/login", GithubLogin)
	router.GET("/project/:id", GetProject)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Index is the Hugoku home page handler will redirect a non logged user to do the loging with Github
// or show a list of projectst and a form to add a project to a logged user,
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Fprint(w, "Hugoku home page!\n")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing the home page template")
	}
	t.Execute(w, nil)
}

func GithubLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "GithubLogin\n")
}

// GetProject is the Hugoku project page handdler and shows the project and the build history.
func GetProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	// TODO: sanitize id
	fmt.Fprintf(w, "GetProject %s!\n", id)
}
