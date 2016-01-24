package routes

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"

	"github.com/gophergala2016/hugoku/ci"
	"github.com/gophergala2016/hugoku/store"
	"github.com/gophergala2016/hugoku/util/cmd"
	"github.com/gophergala2016/hugoku/util/repo"
	"github.com/gophergala2016/hugoku/util/session"
)

// PostProject ...
func PostProject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := session.GetUser(r)
	if err != nil {
		log.Fatal(err)
	}
	projectName := r.PostFormValue("name")
	if projectName == "" {
		log.Println("Missing projectName")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	repoDescription := "Repository for " + projectName + " created by Hugoku"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token.AccessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	log.Printf("Creating %s...", projectName)

	buildInfo, err := ci.Deploy(user.Username, projectName)

	log.Printf("Build duration: %s\n", buildInfo.BuildDuration)

	if err != nil {
		log.Fatalf("Error while trying to create project: %s", err)
	}

	if !repo.Exists(client, user.Username, projectName) {
		gitHubStartTime := time.Now()
		repo := &github.Repository{
			Name:        github.String(projectName),
			Private:     github.Bool(false),
			Description: github.String(repoDescription),
		}
		repo, _, err = client.Repositories.Create("", repo)
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
		err = os.Chdir(wd)
		if err != nil {
			log.Fatal(err)
		}
		githubTime := time.Since(gitHubStartTime)
		log.Printf("Git repo creation duration: %s\n", githubTime)
	}
	project := store.Project{Name: projectName, Description: repoDescription, BuildsInfo: []store.BuildInfo{buildInfo}, LastBuildInfo: buildInfo}
	user.Projects = append(user.Projects, project)

	err = store.SaveUser(user)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetProject is the Hugoku project page handler and shows the project and the build history.
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
	context := struct {
		*store.Project
		Username      string
		GithubProfile string
	}{
		project,
		user.Username,
		user.GithubProfile,
	}
	err = t.Execute(w, context)

	if err != nil {
		log.Fatal("Error executing the project page template:", err)
	}
}

// DeleteProject Deletes a project
func DeleteProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id = ps.ByName("id")
	user, err := session.GetUser(r)
	if err != nil {
		log.Println("Error getting the user from the session: ", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if id == "" {
		log.Println("Missing project id!")
		http.NotFound(w, r)
		return
	}
	_, err = user.GetProject(id)
	if err != nil {
		log.Println("Error getting project: ", id, err)
		http.NotFound(w, r)
	} else {
		// TODO: sanitize id
		wd, _ := os.Getwd()
		err := os.Chdir(wd + "/repos/" + user.Username + "/")
		if err != nil {
			http.Error(w, "Unexpected Error", http.StatusInternalServerError)
			return
		}
		cmd.Run("rm", []string{"-rf", id})

		user.RemoveProject(id)

		err = os.Chdir(wd)
		err = store.SaveUser(user)
		log.Printf("DeleteProject %s!\n", id)
		w.WriteHeader(http.StatusNoContent)
	}
}
