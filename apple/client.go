package apple

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL    = "https://api.music.apple.com"
	ProfileURL = "https://api.music.apple.com/v1/me"
)

type Transport struct {
	DevToken string

	// Base is the base RoundTripper used to make HTTP requests.
	// If nil, http.DefaultTransport is used.
	Base http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with an
// access token from Transport's Source.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	if t.DevToken == "" {
		return nil, errors.New("apple developer token is empty")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.DevToken))

	// req.Body is assumed to be closed by the base RoundTripper.
	reqBodyClosed = true
	return t.base().RoundTrip(req)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

type Client struct {
	DevToken   string
	httpClient *http.Client
}

func NewClient(devToken string) *Client {
	return &Client{
		DevToken: devToken,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: &Transport{DevToken: devToken},
		},
	}
}
