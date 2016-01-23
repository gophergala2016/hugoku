package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"golang.org/x/oauth2"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/gophergala2016/hugoku/store"
)

// random string for oauth2 API calls to protect against CSRF
// TODO: make it random
const OAUTH_RANDOM_CSRF_STRING = "FenaeTaini5thu5eimohpeer1ear5m"

func main() {
	Serve()
}

var (
	// Instantiating the OAuth2 package to exchange the Code for a Token
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("HUGOKU_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("HUGOKU_OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("HUGOKU_AUTH2_CALLBACK_URL"),
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
)

// Serve set the route handlers and serve
func Serve() {

	// Setup middleware
	middle := negroni.Classic()
	store := cookiestore.New([]byte("keyboard cat"))
	middle.Use(sessions.Sessions("hugoku", store))

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/auth/login", githubLoginHandler)
	router.GET("/auth/logout", githubLogoutHandler)
	router.GET("/auth/callback", githubCallbackHandler)
	router.GET("/project/:id", getProjectHandler)
	router.POST("/project", postProjectHandler)
	router.GET("/about", About)
	router.GET("/faq", FAQ)

	// Apply middleware to the router
	middle.UseHandler(router)

	log.Println("Started running on http://127.0.0.1:8080")
	// TODO: Get the port from flag, config file or environment var
	log.Fatal(http.ListenAndServe(":8080", middle))
}

// Index is the Hugoku home page handler will redirect a non logged user to do the loging with Github
// or show a list of projectst and a form to add a project to a logged user,
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	session := sessions.GetSession(r)
	username := session.Get("username")

	log.Println("------", username)

	if username == nil {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal("Error parsing the index page template")
		}
		t.Execute(w, nil)
	} else {
		t, err := template.ParseFiles("templates/home.html")
		if err != nil {
			log.Fatal("Error parsing the home page template")
		}
		t.Execute(w, nil)
	}
}

// githubLoginHandler redirects the user to github to handle the Oauth2 authentication
func githubLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("githubLoginHandler")

	url := oauthConf.AuthCodeURL(OAUTH_RANDOM_CSRF_STRING, oauth2.AccessTypeOnline)
	log.Println("Redirecting the user to github login")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func githubLogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessions.GetSession(r)
	session.Delete("username")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

//gitHubCalbackHandler Called by github after authorization is granted
func githubCallbackHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("githubCallbackHandler")
	state := r.FormValue("state")
	if state != OAUTH_RANDOM_CSRF_STRING {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", OAUTH_RANDOM_CSRF_STRING, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	// Convert the auth code into a token
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get a client that uses the auth token
	oauthClient := oauthConf.Client(oauth2.NoContext, token)
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
	err = store.SaveUser(u)

	session := sessions.GetSession(r)
	session.Set("username", *user.Login)
	session.Set("token", *token)

	log.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// About shows info about the project, team  etc ...
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/about.html")
	if err != nil {
		log.Fatal("Error parsing the about page template")
	}
	t.Execute(w, nil)
}

// FAQ shows frequently asqued questions about the project, team  etc ...
func FAQ(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("templates/faq.html")
	if err != nil {
		log.Fatal("Error parsing the FAQ page template")
	}
	t.Execute(w, nil)
}

// postProjectHandler handles the creation of a new Hugoku project
func postProjectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessions.GetSession(r)
	username := session.Get("username")
	token := session.Get("token")
	projectName := r.PostFormValue("name")
	log.Printf("username: %s, token: %s, projectName: %s\n", username, token, projectName)
	// TODO: Handle the project creation

}

// getProjectHandler is the Hugoku project page handdler and shows the project and the build history.
func getProjectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	// TODO: sanitize id
	log.Printf("getProjectHandler %s!\n", id)
}
