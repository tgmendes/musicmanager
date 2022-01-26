package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/tgmendes/spotistats/auth"
	"github.com/tgmendes/spotistats/repo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	// load environment variables
	pgURL := os.Getenv("POSTGRES_URL")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURL := os.Getenv("SPOTIFY_AUTH_REDIRECT_URL")

	fmt.Println(pgURL)
	conn, err := pgx.Connect(ctx, pgURL)
	if err != nil {
		log.Fatalf("couldn't start PGX: %s\n", err)
	}

	a := auth.NewAuth(clientID, clientSecret, redirectURL, auth.AllScopes())

	h, err := auth.NewHandler(a, &repo.Store{DB: conn})
	if err != nil {
		log.Fatalf("problem starting handler: %s\n", err)
	}

	r := chi.NewRouter()
	r.Get("/authorise", h.Authorise)
	r.Get("/callback", h.AuthCallback)

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
