package handler

import (
	"github.com/tgmendes/soundfuse/apple"
	"github.com/tgmendes/soundfuse/auth"
	"github.com/tgmendes/soundfuse/repo"
)

type Handler struct {
	Store       *repo.Store
	AppleClient *apple.Client
	AppleAuth   *auth.Apple
	SpotifyAuth *auth.Spotify
}
