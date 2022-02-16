package handler

import (
	"fmt"
	"github.com/tgmendes/soundfuse/apple"
	"github.com/tgmendes/soundfuse/auth"
	"github.com/tgmendes/soundfuse/spotify"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type MigrateRequest struct {
	PlaylistID   string `json:"playlist_id"`
	PlaylistHref string `json:"playlist_href"`
	PlaylistName string `json:"playlist_name"`
}

func (h Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "index", nil)
}

func (h Handler) PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	userTokens, ok := auth.TokenFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}
	spotTkn := oauth2.Token{
		AccessToken:  userTokens.SpotifyAccessToken,
		RefreshToken: userTokens.SpotifyRefreshToken,
		Expiry:       time.Now(),
	}
	spotifyClient := spotify.NewClient(h.SpotifyAuth.Client(r.Context(), &spotTkn))
	playlists, err := spotifyClient.GetUsersPlaylists(r.Context(), userTokens.SpotifyUserID, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	applePlaylists, err := h.AppleClient.FetchUserPlaylists(r.Context(), userTokens.AppleUserToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// for idx, playlist := range applePlaylists.Data {
	// 	enrichedPlaylists, err := h.AppleClient.GetPlaylistByPath(r.Context(), cookie.Value, playlist.Href)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(err.Error()))
	// 		return
	// 	}
	// 	applePlaylists.Data[idx].Relationships = enrichedPlaylists.Data[0].Relationships
	// }

	tmplData := struct {
		SpotifyPlaylists spotify.PlaylistsResponse
		ApplePlaylists   apple.LibraryPlaylistsResponse
	}{
		SpotifyPlaylists: playlists,
		ApplePlaylists:   applePlaylists,
	}
	h.renderTemplate(w, "playlists", tmplData)
}

func (h Handler) renderTemplate(w http.ResponseWriter, pagename string, data interface{}) {
	base := filepath.Join("static/templates", "base.html")
	navbar := filepath.Join("static/templates", "navbar.html")
	currPage := filepath.Join("static/templates", fmt.Sprintf("%s.html", pagename))

	tmpl, err := template.ParseFiles(base, navbar, currPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
