package apple

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type StorefrontResponse struct {
	Data []Storefront `json:"data"`
}

type Storefront struct {
	Href string `json:"href"`
	ID   string `json:"id"`
}

func (c Client) GetUserStorefrontID(ctx context.Context, musicUserTkn string) (string, error) {
	url := fmt.Sprintf("%s/storefront", ProfileURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Music-User-Token", musicUserTkn)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code fetching storefront: %d", resp.StatusCode)
	}

	var storeResp StorefrontResponse
	if err := json.NewDecoder(resp.Body).Decode(&storeResp); err != nil {
		return "", fmt.Errorf("unable to unmarshal response: %w", err)
	}

	if len(storeResp.Data) == 0 {
		return "", errors.New("no items in storefront")
	}
	return storeResp.Data[0].ID, nil

}
