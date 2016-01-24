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

	"github.com/gophergala2016/hugoku/ci"
	"github.com/gophergala2016/hugoku/store"
)

// Project is the type of a site
type Project struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

// OAuthRandomCSRString random string for oauth2 API calls to protect against CSRF
// TODO: make it random
const OAuthRandomCSRString = "FenaeTaini5thu5eimohpeer1ear5m"

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
	//router.GET("/v1/projects/:id", getProjectHandler)
	//router.POST("/v1/projects", postProjectHandler)
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
		log.Println("Recovering user data")
		user, err := store.GetUser(username.(string))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(user)
		t.Execute(w, user)
	}
}

// githubLoginHandler redirects the user to github to handle the Oauth2 authentication
func githubLoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("githubLoginHandler")

	url := oauthConf.AuthCodeURL(OAuthRandomCSRString, oauth2.AccessTypeOnline)
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
	if state != OAuthRandomCSRString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", OAuthRandomCSRString, state)
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
	u.Token = *token
	u.AvatarURL = *user.AvatarURL
	err = store.SaveUser(u)

	session := sessions.GetSession(r)
	session.Set("username", *user.Login)

	log.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func postProjectHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	session := sessions.GetSession(r)
	username := session.Get("username")
	user, err := store.GetUser(username.(string))
	var buildStatus = "ok"
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostFormValue("name")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	log.Printf("Creating %s...", projectName)

	_, err = ci.Deploy(username.(string), projectName)

	if err != nil {
		log.Fatalf("Error while trying to create project: %s", err)
		buildStatus = "fail"
	}

	project := store.Project{Name: projectName, BuildInfo: store.BuildInfo{BuildTime: time.Now(), BuildStatus: buildStatus}}
	user.Projects = append(user.Projects, project)

	err = store.SaveUser(user)

	if !repoExists(client, username.(string), projectName) {
		repo := &github.Repository{
			Name:    github.String(projectName),
			Private: github.Bool(false),
		}
		_, _, err = client.Repositories.Create("", repo)

		if err != nil {
			log.Fatalf("Error while trying to create repo: %s", err)
		}
	}

	// TODO: Make git repo to push after
}

// repoExists checks if a repo exists
func repoExists(client *github.Client, username string, projectName string) bool {
	repos, _, _ := client.Repositories.List("", nil)
	projectRepo := username + "/" + projectName
	log.Println(projectRepo)
	for _, b := range repos {
		if *b.FullName == projectRepo {
			log.Println("Github repo already exists, skipping...")
			return true
		}
	}
	return false
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

// getProjectHandler is the Hugoku project page handdler and shows the project and the build history.
func getProjectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	// TODO: sanitize id
	log.Printf("getProjectHandler %s!\n", id)

	if id == "" {
		log.Println("Missing project id!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	t, err := template.ParseFiles("templates/project.html")
	if err != nil {
		log.Fatal("Error parsing the project page template")
	}
	t.Execute(w, nil)

}
