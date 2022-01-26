package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/tgmendes/spotistats/repo"
	"github.com/tgmendes/spotistats/spotify"
	"time"
)

type Playlist struct {
	Store         *repo.Store
	SpotifyClient *spotify.Client
}

func (h *Playlist) CreateMonthlyTopPlaylist(ctx context.Context, userID string) error {
	now := time.Now().UTC()
	internalID := fmt.Sprintf("%s:%s:%d:%s",
		repo.PlaylistTypeShortTerm, now.Month(), now.Year(), userID)

	_, err := h.Store.GetPlaylist(ctx, repo.FilterTypeInternalID, internalID)
	if err == nil {
		return errors.New("playlist alerady created")
	}
	if err != nil && !errors.Is(err, repo.ErrNoResults) {
		return err
	}
	plName := fmt.Sprintf("Top Tracks %s %d", now.Month(), now.Year())
	reqPlaylist := spotify.CreatePlaylistRequest{
		Name:          plName,
		Description:   fmt.Sprintf("The top tracks of %s %d", now.Month(), now.Year()),
		Public:        true,
		Collaborative: false,
	}
	playlist, err := h.SpotifyClient.CreatePlaylist(ctx, userID, reqPlaylist)
	if err != nil {
		return err
	}

	tracks, err := h.SpotifyClient.UserTopTracks(ctx, 50, 0, spotify.ShortTerm)
	if err != nil {
		return err
	}

	var tracksToAdd []string
	for _, track := range tracks.Items {
		tracksToAdd = append(tracksToAdd, track.URI)
	}

	addReq := spotify.AddItemsRequest{
		Position: 0,
		URIs:     tracksToAdd,
	}

	err = h.SpotifyClient.AddItemsToPlaylist(ctx, playlist.ID, addReq)
	if err != nil {
		return err
	}

	newPlaylist := repo.Playlist{
		Name:       plName,
		InternalID: internalID,
		Created:    now,
		Type:       repo.PlaylistTypeShortTerm,
		SpotifyID:  playlist.ID,
		SpotifyURL: playlist.URI,
	}
	err = h.Store.CreatePlaylist(ctx, userID, newPlaylist)
	if err != nil {
		return err
	}
	return nil

}
