package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Track struct {
	Id         string   `json:"id"`
	Genres     []string `json:"genres"`
	Images     []Image  `json:"images"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
	URI        string   `json:"uri"`
}

type TracksResponse struct {
	Items    []Track `json:"items"`
	Previous string  `json:"previous"`
	Next     string  `json:"next"`
	Limit    int     `json:"limit"`
	Offset   int     `json:"offset"`
	Total    int     `json:"total"`
}

func (c *Client) UserTopTracks(ctx context.Context, limit, offset int, timeRange TopTimeRange) (TracksResponse, error) {
	url := fmt.Sprintf("%s/me/top/tracks", BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TracksResponse{}, err
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
		return TracksResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return TracksResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tracks TracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return TracksResponse{}, errors.New("unable to unmarshal response")
	}
	return tracks, nil
}
