package repo

import (
	"context"
	"fmt"
	"time"
)

type Album struct {
	Name            string
	Popularity      int
	ReleaseDate     *time.Time
	TotalTracks     int
	ImageURL        string
	SpotifyID       string
	SpotifyURL      string
	ArtistSpotifyID string
}

func (s *Store) CreateOrUpdateAlbum(ctx context.Context, album Album) error {
	q := `
INSERT INTO albums(name, popularity, release_date, total_tracks, image_url, spotify_id, spotify_url, artist_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, (SELECT artist_id FROM artists WHERE spotify_id = $8))
ON CONFLICT (spotify_id) DO UPDATE SET popularity = EXCLUDED.popularity;

`
	_, err := s.DB.Exec(ctx, q, album.Name, album.Popularity, album.ReleaseDate,
		album.TotalTracks, album.ImageURL, album.SpotifyID, album.SpotifyURL, album.ArtistSpotifyID)
	if err != nil {
		return fmt.Errorf("unable to create album: %w", err)
	}
	return nil
}
