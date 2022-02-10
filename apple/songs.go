package apple

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const ISRCLimit = 25

type TrackResponse struct {
	Href string       `json:"href"`
	Data []TrackData  `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type TrackData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Href       string          `json:"href"`
	Attributes TrackAttributes `json:"attributes"`
}

type TrackAttributes struct {
	Artwork          Artwork  `json:"artwork"`
	ArtistName       string   `json:"artistName"`
	DiscNumber       int      `json:"discNumber"`
	GenreNames       []string `json:"genreNames"`
	DurationInMillis int      `json:"durationInMillis"`
	ReleaseDate      string   `json:"releaseDate"`
	Name             string   `json:"name"`
	ISRC             string   `json:"isrc"`
	HasLyrics        bool     `json:"hasLyrics"`
	AlbumName        string   `json:"albumName"`
}

type SongMeta struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func (c Client) FetchSongsByISRCs(ctx context.Context, isrcs []string) (TrackResponse, error) {
	url := fmt.Sprintf("%s/v1/catalog/gb/songs", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResponse{}, err
	}
	q := req.URL.Query()
	q.Add("filter[isrc]", strings.Join(isrcs, ","))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TrackResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return TrackResponse{}, fmt.Errorf("unexpected status code fetching playlist: %d", resp.StatusCode)
	}

	var trackResp TrackResponse
	if err := json.NewDecoder(resp.Body).Decode(&trackResp); err != nil {
		return TrackResponse{}, fmt.Errorf("unable to unmarshal response: %w", err)
	}
	return trackResp, nil
}
