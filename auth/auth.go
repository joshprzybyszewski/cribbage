package auth

import (
	"context"
	"errors"
	"os"

	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrInvalidService = errors.New(`invalid auth service`)
)

type Service string

const (
	Google Service = `google`
)

// Service provides methods to use OAuth2
type Provider struct {
	services         map[Service]*oauth2.Config
	oauthStateString string
}

func New() *Provider {
	return &Provider{
		services: map[Service]*oauth2.Config{
			Google: {
				RedirectURL:  `http://localhost:8080/auth/` + string(Google) + `/cb`,
				ClientID:     os.Getenv(`OAUTH_GOOGLE_CLIENT_ID`),
				ClientSecret: os.Getenv(`OAUTH_GOOGLE_CLIENT_SECRET`),
				Scopes:       []string{`https://www.googleapis.com/auth/userinfo.email`},
				Endpoint:     google.Endpoint,
			},
		},
		oauthStateString: rand.String(16),
	}
}

func (p *Provider) GetAuthURL(service Service) (string, error) {
	svc, ok := p.services[service]
	if !ok {
		return ``, ErrInvalidService
	}
	return svc.AuthCodeURL(p.oauthStateString), nil
}

func (p *Provider) GetAccessToken(ctx context.Context, service Service, state, code string) (string, error) {
	if state != p.oauthStateString {
		return ``, errors.New(`invalid state string`)
	}
	svc, ok := p.services[service]
	if !ok {
		return ``, ErrInvalidService
	}
	tok, err := svc.Exchange(ctx, code)
	if err != nil {
		return ``, err
	}
	return tok.AccessToken, nil
}
