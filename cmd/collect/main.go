package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/tgmendes/musicmanager/auth"
	"github.com/tgmendes/musicmanager/handler"
	"github.com/tgmendes/musicmanager/repo"
	"github.com/tgmendes/musicmanager/spotify"
	"golang.org/x/oauth2"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	// load environment variables
	pgURL := os.Getenv("POSTGRES_URL")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURL := os.Getenv("SPOTIFY_AUTH_REDIRECT_URL")

	conn, err := pgx.Connect(ctx, pgURL)
	if err != nil {
		panic("no conn")
	}

	a := auth.NewAuth(clientID, clientSecret, redirectURL, auth.AllScopes())

	store := repo.Store{DB: conn}
	tkns, err := store.FetchAll(ctx)
	if err != nil {
		panic(err)
	}

	var token oauth2.Token
	var userID string
	for uID, tkn := range tkns {
		userID = uID
		token = oauth2.Token{
			AccessToken:  tkn.AccessToken,
			RefreshToken: tkn.RefreshToken,
			Expiry:       time.Now(),
		}

	}

	cl := a.Client(ctx, &token)
	spotCl := spotify.NewClient(cl)

	h := handler.Handler{
		Store:         &store,
		SpotifyClient: spotCl,
	}

	err = h.StorePlaylistData(ctx, userID)
	if err != nil {
		panic(err)
	}
}
