package handler

import (
	"context"
	"github.com/tgmendes/soundfuse/repo"
	"github.com/tgmendes/soundfuse/spotify"
)

func (h *Handler) StoreTrack(ctx context.Context, track spotify.Track) error {
	artist := repo.Artist{
		Name:       track.Artists[0].Name,
		SpotifyID:  &track.Artists[0].ID,
		SpotifyURL: &track.Artists[0].Href,
	}

	err := h.Store.CreateOrUpdateArtist(ctx, artist)
	if err != nil {
		return err
	}

	album := repo.Album{
		Name:            track.Album.Name,
		TotalTracks:     track.Album.TotalTracks,
		SpotifyID:       &track.Album.ID,
		SpotifyURL:      &track.Album.HRef,
		ArtistSpotifyID: track.Album.Artists[0].ID,
	}

	err = h.Store.CreateOrUpdateAlbum(ctx, album)
	if err != nil {
		return err
	}

	storeTrack := repo.Track{
		Name:            track.Name,
		Duration:        track.Duration,
		ISRC:            &track.ExternalIDs.ISRC,
		SpotifyID:       &track.ID,
		SpotifyURL:      &track.HRef,
		AlbumSpotifyID:  track.Album.ID,
		ArtistSpotifyID: track.Artists[0].ID,
	}
	err = h.Store.CreateOrUpdateTrack(ctx, storeTrack)
	if err != nil {
		return err
	}
	return nil
}
