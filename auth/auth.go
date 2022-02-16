package auth

import (
	"context"
)

const (
	authState = "spotistats123"
)

type CombinedTokens struct {
	AppleUserToken      string
	SpotifyAccessToken  string
	SpotifyRefreshToken string
	SpotifyUserID       string
}

// TokenFromContext gets the auth token from the context.
func TokenFromContext(ctx context.Context) (CombinedTokens, bool) {
	tokens, ok := ctx.Value("user-tokens").(CombinedTokens)
	return tokens, ok
}
