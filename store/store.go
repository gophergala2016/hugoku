package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
)

// Project type
type Project struct {
	Name       string
	Icon       string
	BuildsInfo []BuildInfo
}

// User type
type User struct {
	Username string
	// Email     string
	Token     oauth2.Token
	AvatarURL string
	Projects  []Project
}

// BuildInfo type
type BuildInfo struct {
	BuildTime   time.Time
	BuildStatus string
}

// GetUser returns an user info
func GetUser(username string) (user User, err error) {
	buf, err := ioutil.ReadFile("./data/" + username + ".json")
	err = json.Unmarshal(buf, &user)
	if err != nil {
		return user, err
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
