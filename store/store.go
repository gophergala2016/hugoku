package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
)

// Project type
type Project struct {
	Name string
}

// User type
type User struct {
	Username string
	Token    oauth2.Token
	Projects []Project
}

// GetUser returns an user info
func GetUser(username string) (user User, err error) {
	buf, err := ioutil.ReadFile("./data/" + user.Username + ".json")
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
