package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Project type
type Project struct {
	Name string
}

// User type
type User struct {
	Username string
	Projects []Project
}

// GetUser returns an user info
func GetUser() (user User, err error) {
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
