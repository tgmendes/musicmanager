package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type User struct {
	ID string `json:"id"`
}

func (c *Client) UserInfo(ctx context.Context) (User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ProfileURL, nil)
	if err != nil {
		return User{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return User{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, errors.New("unable to unmarshal response")
	}
	return user, nil
}
