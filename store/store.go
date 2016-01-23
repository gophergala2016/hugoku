package store

type Repo struct {
	Name string
}

type User struct {
	Username string
	Repos    []Repo
}

// GetAll return all repos from an user
func GetUser() (user User, error err) {
	return user, err
}
