package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"

	"github.com/gophergala2016/hugoku/routes"
)

func main() {
	Serve()
}

// Serve set the route handlers and serve
func Serve() {
	// Setup middleware
	middle := negroni.Classic()
	store := cookiestore.New([]byte("keyboard cat"))
	middle.Use(sessions.Sessions("hugoku", store))

	router := httprouter.New()
	router.GET("/", routes.Index)
	router.GET("/auth/login", routes.GithubLogin)
	router.GET("/auth/logout", routes.GithubLogout)
	router.GET("/auth/callback", routes.GithubLoginCallback)
	router.POST("/project/:id/build", routes.BuildProject)
	router.GET("/project/:id", routes.GetProject)
	router.POST("/project", routes.PostProject)
	router.DELETE("/project/:id", routes.DeleteProject)
	router.GET("/about", routes.About)
	router.GET("/faq", routes.FAQ)

	// Apply middleware to the router
	middle.UseHandler(router)

	log.Println("Started running on http://127.0.0.1:8080")

	// TODO: Get the port from flag, config file or environment var
	log.Fatal(http.ListenAndServe(":8080", middle))
}
