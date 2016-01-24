package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// Project type
type Project struct {
	Name          string
	Icon          string
	Description   string
	LastBuildInfo BuildInfo
	BuildsInfo    []BuildInfo
}

// User type
type User struct {
	Username string
	// Email     string
	Token         oauth2.Token
	AvatarURL     string
	GithubProfile string
	Projects      []Project
}

// BuildInfo type
type BuildInfo struct {
	BuildTime     time.Time
	BuildDuration time.Duration
	BuildStatus   string
	BuildLog      string
	BuildErrorLog string
	BuildPath     string
}

// GetUser returns an user info
func GetUser(username string) (user User, err error) {
	path := "./data/" + username + ".json"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		buf, err := ioutil.ReadFile(path)
		err = json.Unmarshal(buf, &user)
		if err != nil {
			return user, err
		}
	}

	return user, err
}

// SaveUser saves user info
func SaveUser(user User) error {
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile("./data/"+user.Username+".json", b, 0644)
	return err
}

// GetProject gets a project from the user by it's name
func (u User) GetProject(name string) (*Project, error) {
	for _, p := range u.Projects {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, errors.New("Project not found")
}
