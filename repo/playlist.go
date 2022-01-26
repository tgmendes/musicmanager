package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type PlaylistType string

const (
	PlaylistTypeShortTerm  PlaylistType = "short_term"
	PlaylistTypeMediumTerm PlaylistType = "medium_term"
	PlaylistTypeLongTerm   PlaylistType = "long_term"
	PlaylistTypeGeneric    PlaylistType = "long_term"
)

type FilterType string

const (
	FilterTypeInternalID FilterType = "internal_identifier"
	FilterTypeSpotifyID  FilterType = "spotify_id"
)

var ErrNoResults = errors.New("no results found in store")

type Playlist struct {
	Name       string
	InternalID string
	Created    time.Time
	Type       PlaylistType
	SpotifyID  string
	SpotifyURL string
}

func (s *Store) CreatePlaylist(ctx context.Context, userID string, playlist Playlist) error {
	q := `
INSERT INTO playlists(name, internal_identifier, playlist_type, created_date, spotify_id, spotify_url, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);
`
	_, err := s.DB.Exec(ctx, q, playlist.Name, playlist.InternalID, playlist.Type, playlist.Created,
		playlist.SpotifyID, playlist.SpotifyURL, userID)
	if err != nil {
		return fmt.Errorf("unable to create playlist: %w", err)
	}
	return nil
}

func (s *Store) GetPlaylist(ctx context.Context, filterType FilterType, filter string) (Playlist, error) {
	q := `
SELECT name, internal_identifier, playlist_type, created_date, spotify_id, spotify_url
FROM playlists
WHERE %s = $1
`
	q = fmt.Sprintf(q, filterType)
	var playlist Playlist
	row := s.DB.QueryRow(ctx, q, filter)
	err := row.Scan(&playlist.Name, &playlist.InternalID, &playlist.Type, &playlist.Created, &playlist.SpotifyID, &playlist.SpotifyURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Playlist{}, ErrNoResults
		}
		return Playlist{}, fmt.Errorf("unable to scan row: %w", err)
	}

	return playlist, nil
}
