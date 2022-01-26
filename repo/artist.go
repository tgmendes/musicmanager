package repo

import (
	"context"
	"fmt"
)

type Artist struct {
	Name       string
	Popularity int
	ImageURL   string
	SpotifyID  string
	SpotifyURL string
}

func (s *Store) CreateOrUpdateArtist(ctx context.Context, artist Artist) error {
	q := `
INSERT INTO artists(name, popularity, image_url, spotify_id, spotify_url)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (spotify_id) DO UPDATE SET popularity = EXCLUDED.popularity;
`
	_, err := s.DB.Exec(ctx, q, artist.Name, artist.Popularity, artist.ImageURL, artist.SpotifyID, artist.SpotifyURL)
	if err != nil {
		return fmt.Errorf("unable to create artist: %w", err)
	}
	return nil
}
