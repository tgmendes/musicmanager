package handler

import (
	"github.com/tgmendes/soundfuse/apple"
	"github.com/tgmendes/soundfuse/auth"
	"github.com/tgmendes/soundfuse/repo"
	"github.com/tgmendes/soundfuse/worker"
)

type Handler struct {
	Store       *repo.Store
	Cache       *repo.Cache
	AppleClient *apple.Client
	AppleAuth   *auth.Apple
	SpotifyAuth *auth.Spotify
	Worker      *worker.Pool
}
