package repo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type IDMap struct {
	SpotifyID string
	AppleID   string
	AppleType string
}

const RedisExpDur = time.Hour * 24 * 7 // 1 week expiry
type Cache struct {
	Redis *redis.Client
}

func (c *Cache) SetISRCIDs(ctx context.Context, isrc string, ids IDMap) error {
	pipe := c.Redis.Pipeline()

	pipe.HSet(ctx, isrc, "spotify_id", ids.SpotifyID, "apple_id", ids.AppleID, "apple_type", ids.AppleType)
	pipe.Expire(ctx, isrc, RedisExpDur)

	_, err := pipe.Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) ISRCToAppleID(ctx context.Context, isrc string) (IDMap, error) {
	val, err := c.Redis.HGetAll(ctx, isrc).Result()
	if err == redis.Nil {
		return IDMap{}, nil
	}
	if err != nil {
		return IDMap{}, err
	}

	// reset expiry clock
	err = c.Redis.Expire(ctx, isrc, RedisExpDur).Err()
	if err != nil {
		return IDMap{}, err
	}

	return IDMap{
		SpotifyID: val["spotify_id"],
		AppleID:   val["apple_id"],
		AppleType: val["apple_type"],
	}, nil
}
