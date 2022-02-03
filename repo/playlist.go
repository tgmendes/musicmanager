package repo

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"
)

var newPlaylistCols = []string{"name", "playlist_type", "created_date",
	"spotify_id", "spotify_url", "apple_url", "user_id"}

type PlaylistType string

const (
	PlaylistTypeShortTerm  PlaylistType = "short_term"
	PlaylistTypeMediumTerm PlaylistType = "medium_term"
	PlaylistTypeLongTerm   PlaylistType = "long_term"
	PlaylistTypeGeneric    PlaylistType = "generic"
)

type FilterType string

const (
	FilterTypeInternalID FilterType = "internal_identifier"
	FilterTypeSpotifyID  FilterType = "spotify_id"
)

var ErrNoResults = errors.New("no results found in store")

type Playlist struct {
	Name       string
	Created    time.Time
	Type       PlaylistType
	SpotifyID  *string
	SpotifyURL *string
	AppleURL   *string
}

func (s *Store) CreatePlaylist(ctx context.Context, userID string, playlist Playlist) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Insert("playlists").
		Columns(newPlaylistCols...).
		Values(playlist.Name, playlist.Type, playlist.Created, playlist.SpotifyID, playlist.SpotifyURL,
			playlist.AppleURL, userID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to create playlist: %w", err)
	}
	return nil
}

func (s *Store) AddPlaylistTrack(ctx context.Context, playlistSpotifyID, trackSpotifyID string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	playlistQ := playlistBySpotifyIDQuery(playlistSpotifyID)
	trackQ := trackBySpotifyIDQuery(trackSpotifyID)

	sql, args, err := psql.
		Insert("playlist_tracks").
		Columns("playlist_id", "track_id").
		Values(SubQuery(playlistQ), SubQuery(trackQ)).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to create playlist track: %w", err)
	}
	return nil
}

func playlistBySpotifyIDQuery(spotifyID string) sq.SelectBuilder {
	return sq.
		Select("playlist_id").
		From("playlists").
		Where("spotify_id = ?", spotifyID)
}
