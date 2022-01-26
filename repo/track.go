package repo

import (
	"context"
	"fmt"
	"time"
)

type Track struct {
	Name            string
	Popularity      int
	ReleaseDate     *time.Time
	Duration        int
	LastPlayed      *time.Time
	SpotifyID       string
	SpotifyURL      string
	AlbumSpotifyID  string
	ArtistSpotifyID string
}

func (s *Store) CreateOrUpdate(ctx context.Context, track Track) error {
	q := `
INSERT INTO tracks(name, popularity, release_date, duration_ms, last_played, spotify_id, spotify_url, album_id, artist_id)
VALUES (
$1,
$2,
$3,
$4,
$5,
$6,
$7,
(SELECT album_id FROM albums WHERE spotify_id = $8),
(SELECT artist_id FROM artists WHERE spotify_id = $9))
ON CONFLICT (spotify_id) DO UPDATE SET popularity = EXCLUDED.popularity, play_count = tracks.play_count + 1;
`
	_, err := s.DB.Exec(ctx, q, track.Name, track.Popularity, track.ReleaseDate, track.Duration, track.LastPlayed,
		track.SpotifyID, track.SpotifyURL, track.AlbumSpotifyID, track.ArtistSpotifyID)
	if err != nil {
		return fmt.Errorf("unable to create track: %w", err)
	}
	return nil
}
