package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	authState = "spotistats123"
	// AuthURL is the URL to Spotify Accounts Service's OAuth2 endpoint.
	AuthURL = "https://accounts.spotify.com/authorize"
	// TokenURL is the URL to the Spotify Accounts Service's OAuth2
	// token endpoint.
	TokenURL = "https://accounts.spotify.com/api/token"
	// ScopePlaylistReadPrivate seeks permission to read
	// a user's private playlists.
	ScopePlaylistReadPrivate = "playlist-read-private"
	// ScopePlaylistReadCollaborative seeks permission to read
	// a user's collaborative playlists.
	ScopePlaylistReadCollaborative = "playlist-read-collaborative"
	// ScopePlaylistModifyPrivate seeks permission to modify
	// a user's private playlists.
	ScopePlaylistModifyPrivate = "playlist-modify-private"
	// ScopeUserReadPrivate seeks read access to a user's
	// subscription details (type of user account).
	ScopeUserReadPrivate = "user-read-private"
	// ScopeUserReadEmail seeks read access to a user's email address.
	ScopeUserReadEmail = "user-read-email"
	// ScopeUserReadCurrentlyPlaying seeks read access to a user's currently playing track
	ScopeUserReadCurrentlyPlaying = "user-read-currently-playing"
	// ScopeUserReadRecentlyPlayed seeks read access to a user's recently played tracks
	ScopeUserReadRecentlyPlayed = "user-read-recently-played"
	// ScopeUserReadPlaybackState seeks read access to the user's current playback state
	ScopeUserReadPlaybackState = "user-read-playback-state"
	// ScopeUserTopRead seeks read access to a user's top tracks and artists
	ScopeUserTopRead = "user-top-read"
	// ScopeStreaming seeks permission to play music and control playback on your other devices.
	ScopeStreaming = "streaming"
)

func AllScopes() []string {
	return []string{
		ScopePlaylistReadPrivate,
		ScopePlaylistReadCollaborative,
		ScopePlaylistModifyPrivate,
		ScopeUserReadPrivate,
		ScopeUserReadEmail,
		ScopeUserReadCurrentlyPlaying,
		ScopeUserReadRecentlyPlayed,
		ScopeUserReadPlaybackState,
		ScopeUserTopRead,
		ScopeStreaming,
	}
}

type Auth struct {
	config *oauth2.Config
}

func NewAuth(clientID, clientSecret, redirectURL string, scopes []string) *Auth {
	cfg := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      scopes,
	}

	return &Auth{config: &cfg}
}

func (a *Auth) NewToken(ctx context.Context, code, state string) (*oauth2.Token, error) {
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

func (a *Auth) AuthCodeURL() string {
	return a.config.AuthCodeURL(authState)
}

func (a *Auth) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return a.config.Client(ctx, token)
}
