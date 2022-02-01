package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tgmendes/musicmanager/repo"
	"github.com/tgmendes/musicmanager/spotify"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"time"
)

type UserProfile struct {
	ID string `json:"id"`
}

type Handler struct {
	Auth  *Auth
	Store *repo.Store
}

func NewHandler(auth *Auth, store *repo.Store) (*Handler, error) {
	return &Handler{
		Auth:  auth,
		Store: store,
	}, nil
}

func (h *Handler) AuthoriseSpotify(w http.ResponseWriter, _ *http.Request) {
	url := h.Auth.AuthCodeURL()
	t, err := template.ParseFiles("static/authorise.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = t.Execute(w, url)
}

func (h *Handler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	values := r.URL.Query()
	if err := values.Get("error"); err != "" {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	tkn, err := h.Auth.NewToken(ctx, values.Get("code"), values.Get("state"))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}

	err = h.storeToken(r.Context(), tkn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte("token generated"))
}

func (h *Handler) storeToken(ctx context.Context, tkn *oauth2.Token) error {
	httpCl := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spotify.ProfileURL, nil)
	if err != nil {
		return fmt.Errorf("unable to create profile request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tkn.AccessToken)

	resp, err := httpCl.Do(req)
	if err != nil {
		return fmt.Errorf("unable to retrieve profile: %w", err)
	}
	defer resp.Body.Close()
	var prof UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&prof); err != nil {
		return fmt.Errorf("unable to decode profile response: %w", err)
	}

	dbTkn := repo.Token{
		AccessToken:  tkn.AccessToken,
		RefreshToken: tkn.RefreshToken,
	}
	err = h.Store.CreateOrUpdateSpotifyToken(ctx, prof.ID, dbTkn)
	if err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}

	return nil
}
