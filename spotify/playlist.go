package spotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreatePlaylistRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Collaborative bool   `json:"collaborative"`
}

type Playlist struct {
	ID            string         `json:"id"`
	Public        bool           `json:"public"`
	Collaborative bool           `json:"collaborative"`
	Tracks        TracksResponse `json:"tracks"`
	URI           string         `json:"uri"`
}

type AddItemsRequest struct {
	Position int      `json:"position"`
	URIs     []string `json:"uris"`
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
