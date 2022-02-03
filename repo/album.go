package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

var newAlbumCols = []string{"name", "total_tracks", "image_url",
	"spotify_id", "spotify_url", "apple_url", "artist_id"}

type Album struct {
	Name            string
	TotalTracks     int
	ImageURL        *string
	SpotifyID       *string
	SpotifyURL      *string
	AppleURL        *string
	ArtistSpotifyID string
}

func (s *Store) CreateOrUpdateAlbum(ctx context.Context, album Album) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	artistQ := artistBySpotifyIDQuery(album.ArtistSpotifyID)

	sql, args, err := psql.Insert("albums").
		Columns(newAlbumCols...).
		Values(album.Name, album.TotalTracks, album.ImageURL, album.SpotifyID, album.SpotifyURL, album.AppleURL, SubQuery(artistQ)).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to create album: %w", err)
	}

	return nil
}

func albumBySpotifyIDQuery(spotifyID string) sq.SelectBuilder {
	return sq.
		Select("album_id").
		From("albums").
		Where("spotify_id = ?", spotifyID)
}
