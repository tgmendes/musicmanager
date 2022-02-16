package handler

import (
	"fmt"
	"github.com/tgmendes/soundfuse/spotify"
	"net/http"
)

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

	spotCl := spotify.NewClient(h.SpotifyAuth.Client(ctx, tkn))
	user, err := spotCl.UserInfo(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to fetch user: %s", err), http.StatusBadGateway)
	}
	c := http.Cookie{
		Name:  "soundfuse_spotifyaccesstoken",
		Path:  "/",
		Value: tkn.AccessToken,
	}
	http.SetCookie(w, &c)

	c = http.Cookie{
		Name:  "soundfuse_spotifyrefreshtoken",
		Path:  "/",
		Value: tkn.RefreshToken,
	}
	http.SetCookie(w, &c)

	c = http.Cookie{
		Name:  "soundfuse_spotifyuserid",
		Path:  "/",
		Value: user.ID,
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/authorise", 302)
}
