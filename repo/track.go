package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

var newTrackCols = []string{"name", "duration_ms", "isrc",
	"spotify_id", "spotify_url", "apple_url", "album_id", "artist_id"}

type Track struct {
	Name            string
	Duration        int
	ISRC            *string
	SpotifyID       *string
	SpotifyURL      *string
	AppleURL        *string
	AlbumSpotifyID  string
	ArtistSpotifyID string
}

func (s *Store) CreateOrUpdateTrack(ctx context.Context, track Track) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	artistQ := artistBySpotifyIDQuery(track.ArtistSpotifyID)
	albumQ := albumBySpotifyIDQuery(track.AlbumSpotifyID)

	sql, args, err := psql.Insert("tracks").
		Columns(newTrackCols...).
		Values(track.Name, track.Duration, track.ISRC, track.SpotifyID, track.SpotifyURL,
			track.AppleURL, SubQuery(albumQ), SubQuery(artistQ)).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to create track: %w", err)
	}
	return nil
}

func trackByISRCQuery(isrc string) sq.SelectBuilder {
	return sq.
		Select("track_id").
		From("tracks").
		Where("isrc = ?", isrc)
}
