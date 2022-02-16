package handler

import (
	"context"
	"errors"
	"github.com/tgmendes/soundfuse/apple"
	"github.com/tgmendes/soundfuse/auth"
	"github.com/tgmendes/soundfuse/spotify"
	"golang.org/x/oauth2"
	"time"
)

var ErrMissingUserTokens = errors.New("missing user tokens in request context")

type Results struct {
	InCount  int                     `json:"in_count,omitempty"`
	OutCount int                     `json:"out_count,omitempty"`
	Matched  map[string]*SongMatches `json:"matched,omitempty"`
}
type SongMatches struct {
	Name      string `json:"name,omitempty"`
	SpotifyID string `json:"spotify_id,omitempty"`
	AppleID   string `json:"apple_id,omitempty"`
}

type MigrationService struct {
	appleClient   *apple.Client
	spotifyClient *spotify.Client
	userTokens    auth.CombinedTokens

	migrationResults *Results
}

func NewMigrationService(ctx context.Context, spotifyAuth *auth.Spotify, appleClient *apple.Client) (*MigrationService, error) {
	userTokens, ok := auth.TokenFromContext(ctx)
	if !ok {
		return nil, ErrMissingUserTokens
	}

	spotTkn := oauth2.Token{
		AccessToken:  userTokens.SpotifyAccessToken,
		RefreshToken: userTokens.SpotifyRefreshToken,
		Expiry:       time.Now(),
	}
	return &MigrationService{
		appleClient:   appleClient,
		spotifyClient: spotify.NewClient(spotifyAuth.Client(ctx, &spotTkn)),
		userTokens:    userTokens,
	}, nil
}

func (m *MigrationService) Migrate(ctx context.Context, req MigrateRequest) (*Results, error) {
	playlistItems, err := m.spotifyClient.GetPlaylistItems(ctx, req.PlaylistHref)
	if err != nil {
		return nil, err
	}

	storefrontID, err := m.appleClient.GetUserStorefrontID(ctx, m.userTokens.AppleUserToken)
	if err != nil {
		return nil, err
	}

	m.migrationResults = &Results{
		InCount: len(playlistItems.Items),
		Matched: map[string]*SongMatches{},
	}

	tracksToAdd, err := m.spotifyToAppleTracks(ctx, storefrontID, playlistItems.Items)
	if err != nil {
		return nil, err
	}

	m.migrationResults.OutCount = len(tracksToAdd)

	err = m.CreatePlaylist(ctx, req.PlaylistName, tracksToAdd)
	if err != nil {
		return nil, err
	}
	return m.migrationResults, nil
}

func (m *MigrationService) spotifyToAppleTracks(ctx context.Context, storefrontID string, spotifyTracks []spotify.TrackItem) ([]apple.RelationshipTrack, error) {
	var appleTracks []apple.RelationshipTrack
	var isrcs []string
	for i, item := range spotifyTracks {
		m.migrationResults.Matched[item.Track.ExternalIDs.ISRC] = &SongMatches{
			Name:      item.Track.Name,
			SpotifyID: item.Track.ID,
		}
		isrcs = append(isrcs, item.Track.ExternalIDs.ISRC)

		if len(isrcs) == apple.ISRCLimit-1 || i == len(spotifyTracks)-1 {
			songsResp, err := m.appleClient.FetchSongsByISRCs(ctx, storefrontID, isrcs)
			if err != nil {
				return nil, err
			}

			seenISRC := map[string]struct{}{}
			for _, track := range songsResp.Data {
				if _, ok := seenISRC[track.Attributes.ISRC]; ok {
					continue
				}
				if _, ok := m.migrationResults.Matched[track.Attributes.ISRC]; ok {
					m.migrationResults.Matched[track.Attributes.ISRC].AppleID = track.ID
				}
				seenISRC[track.Attributes.ISRC] = struct{}{}
				toAppend := apple.RelationshipTrack{
					ID:   track.ID,
					Type: track.Type,
				}
				appleTracks = append(appleTracks, toAppend)
			}
			isrcs = nil
		}
	}
	return appleTracks, nil
}

func (m *MigrationService) CreatePlaylist(ctx context.Context, name string, tracks []apple.RelationshipTrack) error {
	plReq := apple.CreatePlaylistRequest{
		Attributes: apple.CreatePlaylistAttributes{
			Name:        name,
			Description: "Music manager auto generated playlist from spotify",
		},
		Relationships: apple.CreatePlaylistRelationships{
			Tracks: apple.TracksData{
				Data: tracks,
			},
		},
	}

	err := m.appleClient.CreateUserPlaylist(ctx, m.userTokens.AppleUserToken, plReq)
	if err != nil {
		return err
	}
	return nil
}
