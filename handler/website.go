package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tgmendes/musicmanager/apple"
	"github.com/tgmendes/musicmanager/spotify"
	"html/template"
	"net/http"
	"path/filepath"
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
	cookie, err := r.Cookie("appleToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/authorise", 302)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	playlists, err := h.SpotifyClient.GetUsersPlaylists(r.Context(), "a2mqb93izo81vhk8ijacgz73c", 50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[spotify] %s", err)))
		return
	}

	applePlaylists, err := h.AppleClient.FetchUserPlaylists(r.Context(), cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("[apple] %s", err)))
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

func (h Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	url := h.SpotifyAuth.AuthCodeURL()
	tmplData := struct {
		AppleDevToken  string
		SpotifyAuthURL string
	}{
		AppleDevToken:  h.AppleClient.DevToken,
		SpotifyAuthURL: url,
	}

	h.renderTemplate(w, "authorise", tmplData)
}

func (h *Handler) SpotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	values := r.URL.Query()
	if err := values.Get("error"); err != "" {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}

	tkn, err := h.SpotifyAuth.NewToken(ctx, values.Get("code"), values.Get("state"))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}

	tknString, err := json.Marshal(tkn)
	if err != nil {
		http.Error(w, "unable to generate token string", http.StatusInternalServerError)
	}
	c := http.Cookie{
		Name:  "SpotifyUserToken",
		Value: string(tknString),
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/authorise", 302)
}

func (h Handler) Migrate(w http.ResponseWriter, r *http.Request) {
	appleMusicTkn, err := r.Cookie("appleToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	var req MigrateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	playlist, err := h.SpotifyClient.GetPlaylistItems(r.Context(), req.PlaylistHref)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var tracksToAdd []apple.RelationshipTrack
	var isrcs []string
	var count int
	for _, item := range playlist.Items {
		if count == apple.ISRCLimit-1 {
			tracks, err := h.fetchISRCSongs(r.Context(), isrcs)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			tracksToAdd = append(tracksToAdd, tracks...)

			isrcs = nil
			count = 0
		}

		isrcs = append(isrcs, item.Track.ExternalIDs.ISRC)
		count++
	}

	tracks, err := h.fetchISRCSongs(r.Context(), isrcs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tracksToAdd = append(tracksToAdd, tracks...)
	err = h.createPlaylist(r.Context(), appleMusicTkn.Value, req.PlaylistName, tracksToAdd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) renderTemplate(w http.ResponseWriter, pagename string, data interface{}) {
	base := filepath.Join("static/templates", "base.html")
	navbar := filepath.Join("static/templates", "navbar.html")
	currPage := filepath.Join("static/templates", fmt.Sprintf("%s.html", pagename))

	tmpl, err := template.ParseFiles(base, navbar, currPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (h Handler) fetchISRCSongs(ctx context.Context, isrcs []string) ([]apple.RelationshipTrack, error) {
	songsResp, err := h.AppleClient.FetchSongsByISRCs(ctx, isrcs)
	if err != nil {
		return nil, err
	}

	var tracks []apple.RelationshipTrack
	seenISRC := map[string]struct{}{}
	for _, track := range songsResp.Data {
		if _, ok := seenISRC[track.Attributes.ISRC]; ok {
			continue
		}
		seenISRC[track.Attributes.ISRC] = struct{}{}
		toAppend := apple.RelationshipTrack{
			ID:   track.ID,
			Type: track.Type,
		}
		tracks = append(tracks, toAppend)
	}

	return tracks, nil
}

func (h Handler) createPlaylist(ctx context.Context, userTkn, name string, tracks []apple.RelationshipTrack) error {
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
	err := h.AppleClient.CreateUserPlaylist(ctx, userTkn, plReq)
	if err != nil {
		return err
	}
	return nil
}
