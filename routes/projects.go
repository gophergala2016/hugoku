package routes

import (
	"log"
	"net/http"
	"text/template"
	"os"
	"time"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"

	"github.com/gophergala2016/hugoku/ci"
	"github.com/gophergala2016/hugoku/store"
	"github.com/gophergala2016/hugoku/util/session"
	"github.com/gophergala2016/hugoku/util/cmd"
	"github.com/gophergala2016/hugoku/util/repo"
)

func PostProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var buildStatus = "ok"
	user, err := session.GetUser(r)
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

	startTime := time.Now()
	//path, err := ci.Deploy(username.(string), projectName)
	_, err = ci.Deploy(user.Username, projectName)

	buildDuration := time.Since(startTime)
	if err != nil {
		log.Fatalf("Error while trying to create project: %s", err)
		buildStatus = "fail"
	}

	buildInfo := store.BuildInfo{BuildTime: time.Now(), BuildDuration: buildDuration, BuildStatus: buildStatus}
	project := store.Project{Name: projectName, BuildsInfo: []store.BuildInfo{buildInfo}, LastBuildInfo: buildInfo}
	user.Projects = append(user.Projects, project)

	err = store.SaveUser(user)

	if !repo.Exists(client, user.Username, projectName) {
		repo := &github.Repository{
			Name:    github.String(projectName),
			Private: github.Bool(false),
		}
		_, _, err = client.Repositories.Create("", repo)
		if err != nil {
			log.Fatalf("Error while trying to create repo: %s", err)
		}

		// Push the repo
		wd, _ := os.Getwd()
		err := os.Chdir(wd + "/repos/" + user.Username + "/" + projectName + "/")
		if err != nil {
			log.Fatal(err)
		}
		cmd.Run("git", []string{"init", "--quiet"})
		cmd.Run("git", []string{"add", "."})
		cmd.Run("git", []string{"commit", "-m", "'initial source code'"})
		cmd.Run("git", []string{"remote", "add", "origin", "git@github.com:" + user.Username + "/" + projectName + ".git"})
		cmd.Run("git", []string{"push", "--quiet", "-u", "origin", "master"})
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// getProject is the Hugoku project page handdler and shows the project and the build history.
func GetProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	user, err := session.GetUser(r)
	if err != nil {
		log.Fatal("Error getting the user from the session: ", err)
	}
	// TODO: sanitize id
	log.Printf("getProject %s!\n", id)

	if id == "" {
		log.Println("Missing project id!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	project, err := user.GetProject(id)
	if err != nil {
		log.Println("Error getting project: ", id, err)
	}
	t, err := template.ParseFiles("templates/project.html",
		"templates/partials/footer.html",
		"templates/partials/header.html")
	if err != nil {
		log.Fatal("Error parsing the project page template:", err)
	}
	err = t.Execute(w, project)
	if err != nil {
		log.Fatal("Error executing the project page template:", err)
	}
}
