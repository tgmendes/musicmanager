package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/tgmendes/musicmanager/apple"
	"github.com/tgmendes/musicmanager/auth"
	"github.com/tgmendes/musicmanager/handler"
	"github.com/tgmendes/musicmanager/repo"
	"github.com/tgmendes/musicmanager/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pgURL := os.Getenv("POSTGRES_URL")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURL := os.Getenv("SPOTIFY_AUTH_REDIRECT_URL")
	appleIss := os.Getenv("APPLE_ISSUER")
	appleKID := os.Getenv("APPLE_KID")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, pgURL)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	a := auth.NewAuth(clientID, clientSecret, redirectURL, auth.AllScopes())

	store := repo.Store{DB: conn}
	tkns, err := store.FetchAll(ctx)
	if err != nil {
		log.Fatalf("could not fetch user tokens: %s", err)
	}

	var token oauth2.Token
	for _, tkn := range tkns {
		token = oauth2.Token{
			AccessToken:  tkn.AccessToken,
			RefreshToken: tkn.RefreshToken,
			Expiry:       time.Now(),
		}
	}

	cl := a.Client(ctx, &token)
	spotCl := spotify.NewClient(cl)

	appleTkn := auth.GenerateToken(appleIss, appleKID)

	p8key, err := ioutil.ReadFile("AuthKey_MTY4WUTFNX.p8")
	if err != nil {
		log.Fatalf("unable to open dev token: %s", err)
	}
	signedtkn, err := auth.GenerateSignedToken(appleTkn, p8key)
	if err != nil {
		log.Fatalf("unable to generate signed token: %s", err)
	}

	appleCl := apple.NewClient(signedtkn)
	h := handler.Handler{
		Store:         &store,
		SpotifyClient: spotCl,
		AppleClient:   appleCl,
		SpotifyAuth:   a,
	}

	fs := http.FileServer(http.Dir("./static/"))
	r := chi.NewRouter()
	r.Handle("/static/*", http.StripPrefix("/static", fs))
	r.Get("/", h.IndexHandler)
	r.Get("/authorise", h.AuthHandler)
	r.Get("/playlists", h.PlaylistHandler)
	r.Post("/migrate", h.Migrate)
	r.Get("/callback", h.SpotifyCallbackHandler)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the service listening for api requests.
	go func() {
		log.Println("listening on port :8080")
		serverErrors <- srv.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("server encountered error: %s\n", err)

	case sig := <-shutdown:
		log.Printf("initialising shutdown: %s\n", sig)
		defer log.Printf("shutdown complete: %s\n", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			log.Fatalf("could not stop server gracefully: %s", err)
		}
	}
}
