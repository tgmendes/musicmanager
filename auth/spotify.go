package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/tgmendes/soundfuse/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

type Spotify struct {
	config *oauth2.Config
}

func NewSpotify(clientID, clientSecret, redirectURL string, scopes []string) *Spotify {
	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotify.AuthURL,
			TokenURL: spotify.TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}

	return &Spotify{config: &cfg}
}

func (a *Spotify) NewToken(ctx context.Context, code, state string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("provided code is empty")
	}
	if state != authState {
		return nil, errors.New("spotify: redirect state parameter doesn't match")
	}
	tkn, err := a.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange token: %w", err)
	}

	return tkn, nil
}

func (a *Spotify) AuthCodeURL() string {
	return a.config.AuthCodeURL(authState)
}

func (a *Spotify) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return a.config.Client(ctx, token)
}
