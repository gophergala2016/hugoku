package oauth

import (
	"os"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	// Conf instantiating the OAuth2 package to exchange the Code for a Token
	Conf = &oauth2.Config{
		ClientID:     os.Getenv("HUGOKU_OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("HUGOKU_OAUTH2_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("HUGOKU_OAUTH2_CALLBACK_URL"),
		Scopes:       []string{"user", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
)
