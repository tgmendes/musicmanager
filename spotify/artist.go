package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Artist struct {
	Id         string   `json:"id"`
	Genres     []string `json:"genres"`
	Images     []Image  `json:"images"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
	Uri        string   `json:"uri"`
}

type ArtistsResponse struct {
	Items    []Artist `json:"items"`
	Previous string   `json:"previous"`
	Next     string   `json:"next"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Total    int      `json:"total"`
}

func (c *Client) UserTopArtists(ctx context.Context, limit, offset int, timeRange TopTimeRange) (ArtistsResponse, error) {
	url := fmt.Sprintf("%s/me/top/artists", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ArtistsResponse{}, err
	}
	q := req.URL.Query()
	q.Add("time_range", string(timeRange))
	if limit != 0 {
		q.Add("limit", strconv.Itoa(limit))
	}

	if offset != 0 {
		q.Add("offset", strconv.Itoa(offset))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ArtistsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return ArtistsResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var artists ArtistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return ArtistsResponse{}, errors.New("unable to unmarshal response")
	}
	return artists, nil
}
