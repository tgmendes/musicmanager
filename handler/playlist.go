package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/tgmendes/musicmanager/repo"
	"github.com/tgmendes/musicmanager/spotify"
	"time"
)

func (h Handler) StorePlaylistData(ctx context.Context, userID string) error {
	playlistData, err := h.SpotifyClient.GetUsersPlaylists(ctx, userID, 50)
	if err != nil {
		return err
	}

	counter := 1
	for _, playlist := range playlistData.Items {
		newPlaylist := repo.Playlist{
			Name:       playlist.Name,
			Type:       repo.PlaylistTypeGeneric,
			Created:    time.Now().UTC(),
			SpotifyID:  &playlist.ID,
			SpotifyURL: &playlist.Href,
		}
		err := h.Store.CreatePlaylist(ctx, userID, newPlaylist)
		if err != nil {
			return err
		}

		err = h.iterPlaylistTracks(ctx, playlist, h.StoreTrack)
		if err != nil {
			return err
		}
		if counter == 1 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

//
// func (h *Handler) CreateMonthlyTopPlaylist(ctx context.Context, userID string) error {
// 	now := time.Now().UTC()
// 	internalID := fmt.Sprintf("%s:%s:%d:%s",
// 		repo.PlaylistTypeShortTerm, now.Month(), now.Year(), userID)
//
// 	_, err := h.Store.GetPlaylist(ctx, repo.FilterTypeInternalID, internalID)
// 	if err == nil {
// 		return errors.New("playlist alerady created")
// 	}
// 	if err != nil && !errors.Is(err, repo.ErrNoResults) {
// 		return err
// 	}
// 	plName := fmt.Sprintf("Top Tracks %s %d", now.Month(), now.Year())
// 	reqPlaylist := spotify.CreatePlaylistRequest{
// 		Name:          plName,
// 		Description:   fmt.Sprintf("The top tracks of %s %d", now.Month(), now.Year()),
// 		Public:        true,
// 		Collaborative: false,
// 	}
// 	playlist, err := h.SpotifyClient.CreatePlaylist(ctx, userID, reqPlaylist)
// 	if err != nil {
// 		return err
// 	}
//
// 	tracks, err := h.SpotifyClient.UserTopTracks(ctx, 50, 0, spotify.ShortTerm)
// 	if err != nil {
// 		return err
// 	}
//
// 	var tracksToAdd []string
// 	for _, track := range tracks.Items {
// 		tracksToAdd = append(tracksToAdd, track.HRef)
// 	}
//
// 	addReq := spotify.AddItemsRequest{
// 		Position: 0,
// 		URIs:     tracksToAdd,
// 	}
//
// 	err = h.SpotifyClient.AddItemsToPlaylist(ctx, playlist.ID, addReq)
// 	if err != nil {
// 		return err
// 	}
//
// 	newPlaylist := repo.Playlist{
// 		Name:       plName,
// 		InternalID: internalID,
// 		Created:    now,
// 		Type:       repo.PlaylistTypeShortTerm,
// 		SpotifyID:  playlist.ID,
// 		SpotifyURL: playlist.URI,
// 	}
// 	err = h.Store.CreatePlaylist(ctx, userID, newPlaylist)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
//
// }

func (h *Handler) iterPlaylistTracks(ctx context.Context, playlist spotify.Playlist, callbackFn func(ctx context.Context, track spotify.Track) error) error {
	currItems, err := h.SpotifyClient.GetPlaylistItems(ctx, playlist.Items.Href)
	if err != nil {
		return err
	}

	count := 0
	for {
		for _, item := range currItems.Items {
			err := callbackFn(ctx, item.Track)
			if err != nil {
				return err
			}

			err = h.Store.AddPlaylistTrack(ctx, playlist.ID, item.Track.ID)
			if err != nil {
				return err
			}
		}

		if currItems.Next == "" {
			return nil
		}
		currItems, err = h.SpotifyClient.GetPlaylistItems(ctx, currItems.Next)
		if err != nil {
			return err
		}
		count++
		if count > 5 {
			return errors.New("exceeded loop limit")
		}
		fmt.Printf("Next tracks: %s\n", currItems.Next)
	}
}