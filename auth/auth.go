package model

import (
	"golang.org/x/oauth2"
)

type oauthConfig struct {
	google *oauth2.Config
}

type Auth struct {
	config oauthConfig
}
