package middleware

import (
	"context"
	"github.com/tgmendes/soundfuse/auth"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appleCookie, err := r.Cookie("soundfuse_appleusertoken")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/authorise", http.StatusFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spotifyAccCookie, err := r.Cookie("soundfuse_spotifyaccesstoken")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/authorise", 302)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spotifyRefCookie, err := r.Cookie("soundfuse_spotifyrefreshtoken")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/authorise", 302)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		spotifyUserID, err := r.Cookie("soundfuse_spotifyuserid")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/authorise", 302)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tokens := auth.CombinedTokens{
			AppleUserToken:      appleCookie.Value,
			SpotifyAccessToken:  spotifyAccCookie.Value,
			SpotifyRefreshToken: spotifyRefCookie.Value,
			SpotifyUserID:       spotifyUserID.Value,
		}
		ctx := context.WithValue(r.Context(), "user-tokens", tokens)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
