package spotify

import (
	"net/http"
	"time"
)

const (
	BaseURL    = "https://api.spotify.com/v1"
	ProfileURL = "https://api.spotify.com/v1/me"
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
