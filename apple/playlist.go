package apple

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LibraryPlaylistsResponse struct {
	Next string             `json:"next"`
	Data []LibraryPlaylists `json:"data"`
	Meta ResponseMeta       `json:"meta"`
}

type LibraryPlaylists struct {
	ID            string             `json:"id"`
	Href          string             `json:"href"`
	Attributes    PlaylistAttributes `json:"attributes"`
	Relationships Relationships      `json:"relationships"`
}

type PlaylistAttributes struct {
	Name      string    `json:"name"`
	IsPublic  bool      `json:"isPublic"`
	Artwork   Artwork   `json:"artwork"`
	DateAdded time.Time `json:"dateAdded"`
}

type Relationships struct {
	Tracks TrackResponse `json:"tracks"`
}

type CreatePlaylistRequest struct {
	Attributes    CreatePlaylistAttributes    `json:"attributes"`
	Relationships CreatePlaylistRelationships `json:"relationships"`
}

type CreatePlaylistAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreatePlaylistRelationships struct {
	Tracks TracksData `json:"tracks"`
}

type TracksData struct {
	Data []RelationshipTrack `json:"data"`
}
type RelationshipTrack struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ParentData struct {
	Data RelationshipParent `json:"data"`
}

type RelationshipParent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c Client) CreateUserPlaylist(ctx context.Context, musicUserTkn string, plReq CreatePlaylistRequest) error {
	url := fmt.Sprintf("%s/library/playlists", ProfileURL)
	reqBody, err := json.Marshal(plReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Music-User-Token", musicUserTkn)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(bB))
	return nil
}

func (c Client) FetchUserPlaylists(ctx context.Context, musicUserTkn string) (LibraryPlaylistsResponse, error) {
	url := fmt.Sprintf("%s/library/playlists", ProfileURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return LibraryPlaylistsResponse{}, err
	}
	q := req.URL.Query()
	q.Add("extend", "artwork")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Music-User-Token", musicUserTkn)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LibraryPlaylistsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return LibraryPlaylistsResponse{}, fmt.Errorf("unexpected status code fetching playlist: %d", resp.StatusCode)
	}

	var plResp LibraryPlaylistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&plResp); err != nil {
		return LibraryPlaylistsResponse{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return plResp, nil
}

func (c Client) GetPlaylistByPath(ctx context.Context, musicUserTkn, path string) (LibraryPlaylistsResponse, error) {
	url := fmt.Sprintf("%s%s", BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return LibraryPlaylistsResponse{}, err
	}

	q := req.URL.Query()
	q.Add("extend", "artwork")
	q.Add("include", "tracks")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Music-User-Token", musicUserTkn)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LibraryPlaylistsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return LibraryPlaylistsResponse{}, fmt.Errorf("unexpected status code fetching playlist: %d", resp.StatusCode)
	}

	var plResp LibraryPlaylistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&plResp); err != nil {
		return LibraryPlaylistsResponse{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return plResp, nil
}
