package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

var newArtistCols = []string{"name", "image_url", "spotify_id", "spotify_url", "apple_url"}

type Artist struct {
	Name       string
	ImageURL   *string
	SpotifyID  *string
	SpotifyURL *string
	AppleURL   *string
}

func (s *Store) CreateOrUpdateArtist(ctx context.Context, artist Artist) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.
		Insert("artists").
		Columns(newArtistCols...).
		Values(artist.Name, nil, artist.SpotifyID, artist.SpotifyURL, nil).
		Suffix("ON CONFLICT DO NOTHING;").
		ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to create artist: %w", err)
	}
	return nil
}

func artistBySpotifyIDQuery(spotifyID string) sq.SelectBuilder {
	return sq.
		Select("artist_id").
		From("artists").
		Where("spotify_id = ?", spotifyID)
}
