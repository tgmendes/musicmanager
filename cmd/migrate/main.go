package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/tgmendes/musicmanager/auth"
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

	tknRepo := repo.Store{DB: conn}
	tkns, err := tknRepo.FetchAll(ctx)
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
	playlists, err := spotCl.GetUsersPlaylists(ctx, userID, 50)
	if err != nil {
		panic(err)
	}

	tracks, err := spotCl.GetPlaylistItems(ctx, playlists.Items[0].Items.Href)
	if err != nil {
		panic(err)
	}
	for _, item := range tracks.Items {
		fmt.Printf("%s by %s on album %s\n", item.Track.Name, item.Track.Artists[0].Name, item.Track.Album.Name)
		// fmt.Printf("Track Name: %s\n", item.Track.Name)
	}
	// fmt.Printf("Showing %d out of %d playlists.\n", playlists.Limit, playlists.Total)
	// fmt.Printf("Next playlist: %s.\n", playlists.Next)
	// for _, playlist := range playlists.Items {
	// 	fmt.Printf("%s: %d tracks\n", playlist.Name, playlist.Items.Total)
	// 	fmt.Printf("tracks URL: %s\n", playlist.Items.Href)
	// 	tracks, err := spotCl.GetPlaylistItems(ctx, playlist.Items.Href)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("next: %s", tracks.Next)
	// 	for _, item := range tracks.Items {
	// 		fmt.Printf("Track Name: %s\n", item.Track.Name)
	// 	}
	// 	time.Sleep(time.Second)
	// }
}
