package repo

import (
	"context"
	"fmt"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

func (s *Store) CreateOrUpdateSpotifyToken(ctx context.Context, userID string, token Token) error {
	q := `
INSERT INTO spotify_tokens(access_token, refresh_token, user_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE SET access_token = EXCLUDED.access_token;
`
	_, err := s.DB.Exec(ctx, q, token.AccessToken, token.RefreshToken, userID)
	if err != nil {
		return fmt.Errorf("unable to create token: %w", err)
	}
	return nil
}

func (s *Store) FetchAll(ctx context.Context) (map[string]Token, error) {
	q := "SELECT user_id, access_token, refresh_token FROM spotify_tokens LIMIT 100;"
	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch tokens: %w", err)
	}
	defer rows.Close()

	results := map[string]Token{}
	for rows.Next() {
		var tkn Token
		var userID *string
		_ = rows.Scan(&userID, &tkn.AccessToken, &tkn.RefreshToken)
		results[*userID] = tkn
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error scanning rows: %w", rows.Err())
	}
	return results, nil
}
