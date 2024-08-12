package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{client: client}
}

func (uc *UserCache) Set(ctx context.Context, openID string, locationObj Location) error {
	data, err := json.Marshal(locationObj)
	if err != nil {
		return err
	}
	return uc.client.Set(ctx, "user:"+openID, data, 24*time.Hour).Err()
}

func (uc *UserCache) Get(ctx context.Context, openID string) (Location, bool, error) {
	data, err := uc.client.Get(ctx, "user:"+openID).Bytes()
	if err == redis.Nil {
		return Location{}, false, nil
	} else if err != nil {
		return Location{}, false, err
	}
	var location Location
	err = json.Unmarshal(data, &location)
	return location, true, err
}

func (uc *UserCache) GetAll(ctx context.Context) (map[string]Location, error) {
	keys, err := uc.client.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}
	locations := make(map[string]Location)
	for _, key := range keys {
		openID := key[5:] // remove "user:"
		location, _, err := uc.Get(ctx, openID)
		if err != nil {
			return nil, err
		}
		locations[openID] = location
	}
	return locations, nil
}
