package repo

import (
	"log"

	"github.com/google/go-github/github"
)

// Exists checks if a repo exists
func Exists(client *github.Client, username string, projectName string) bool {
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
