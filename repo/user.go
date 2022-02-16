package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

var newUserCols = []string{"spotify_user_id", "apple_storefront_id"}

type User struct {
	ID                int
	SpotifyID         string
	AppleStorefrontID string
}

func (s Store) CreateUser(ctx context.Context, user User) (*User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert("users").
		Columns(newUserCols...).
		Values(user.SpotifyID, user.AppleStorefrontID).
		Suffix("ON CONFLICT DO NOTHING").
		Suffix("RETURNING \"user_id\"").
		ToSql()
	if err != nil {
		return nil, err
	}
	var uID *int
	err = s.DB.QueryRow(ctx, sql, args...).Scan(&uID)
	if err != nil {
		return nil, err
	}

	u := User{
		SpotifyID:         user.SpotifyID,
		AppleStorefrontID: user.AppleStorefrontID}
	if uID != nil {
		u.ID = *uID
	}
	return &u, nil

}
