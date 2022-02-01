package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CreatePlaylistRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Collaborative bool   `json:"collaborative"`
}

type PlaylistsResponse struct {
	Href     string     `json:"href"`
	Items    []Playlist `json:"items"`
	Previous string     `json:"previous"`
	Next     string     `json:"next"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Total    int        `json:"total"`
}

type Playlist struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Public        bool          `json:"public"`
	Collaborative bool          `json:"collaborative"`
	Items         PlaylistItems `json:"tracks"`
	URI           string        `json:"uri"`
}

type PlaylistItems struct {
	Href     string      `json:"href"`
	Items    []TrackItem `json:"items"`
	Previous string      `json:"previous"`
	Next     string      `json:"next"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Total    int         `json:"total"`
}

type AddItemsRequest struct {
	Position int      `json:"position"`
	URIs     []string `json:"uris"`
}

func (c *Client) GetUsersPlaylists(ctx context.Context, userID string, limit int) (PlaylistsResponse, error) {
	url := fmt.Sprintf("%s/users/%s/playlists", BaseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return PlaylistsResponse{}, err
	}
	q := req.URL.Query()
	if limit != 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PlaylistsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PlaylistsResponse{}, fmt.Errorf("unexpected status code fetching playlist: %d", resp.StatusCode)
	}

	var plResp PlaylistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&plResp); err != nil {
		return PlaylistsResponse{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return plResp, nil
}

func (c *Client) CreatePlaylist(ctx context.Context, userID string, reqPlaylist CreatePlaylistRequest) (Playlist, error) {
	url := fmt.Sprintf("%s/users/%s/playlists", BaseURL, userID)
	body, err := json.Marshal(reqPlaylist)
	if err != nil {
		return Playlist{}, fmt.Errorf("unable to marshal body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return Playlist{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Playlist{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return Playlist{}, fmt.Errorf("unexpected status code creating playlist: %d", resp.StatusCode)
	}

	var playlist Playlist
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return Playlist{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return playlist, nil
}

func (c *Client) GetPlaylistItems(ctx context.Context, url string) (PlaylistItems, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return PlaylistItems{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PlaylistItems{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PlaylistItems{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var items PlaylistItems
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return PlaylistItems{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return items, nil
}

func (c *Client) AddItemsToPlaylist(ctx context.Context, playlistID string, addReq AddItemsRequest) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", BaseURL, playlistID)
	body, err := json.Marshal(addReq)
	if err != nil {
		return fmt.Errorf("unable to marshal body: %w", err)
	}
	fmt.Println(string(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return fmt.Errorf("unexpected status code adding items playlist: %d", resp.StatusCode)
	}

	return nil
}
