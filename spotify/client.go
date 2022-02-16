package spotify

import (
	"net/http"
	"time"
)

const (
	BaseURL    = "https://api.spotify.com/v1"
	ProfileURL = "https://api.spotify.com/v1/me"
	// AuthURL is the URL to Spotify Accounts Service's OAuth2 endpoint.
	AuthURL = "https://accounts.spotify.com/authorize"
	// TokenURL is the URL to the Spotify Accounts Service's OAuth2
	// token endpoint.
	TokenURL = "https://accounts.spotify.com/api/token"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Client{httpClient: httpClient}
}
