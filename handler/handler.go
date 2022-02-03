package handler

import (
	"github.com/tgmendes/musicmanager/repo"
	"github.com/tgmendes/musicmanager/spotify"
)

type Handler struct {
	Store         *repo.Store
	SpotifyClient *spotify.Client
}
