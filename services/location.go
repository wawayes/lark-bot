package services

import (
	"context"
	"fmt"
	"time"
)

type Location struct {
	Name       string
	Latitude   string
	Longtitude string
	Timestamp  time.Time
}

type LocationService struct {
	cache *UserCache
}

func NewLocationService(cache *UserCache) *LocationService {
	return &LocationService{
		cache: cache,
	}
}

func (ls *LocationService) SetLocation(ctx context.Context, openID string, locationObj Location) error {
	fmt.Printf("SetLocation: %v\n", locationObj)
	return ls.cache.Set(ctx, openID, locationObj)
}

func (ls *LocationService) GetLocation(ctx context.Context, openID string) (Location, bool) {
	location, exists, _ := ls.cache.Get(ctx, openID)
	return location, exists
}

func (ls *LocationService) GetAllLocations(ctx context.Context) map[string]Location {
	Locations, _ := ls.cache.GetAll(ctx)
	return Locations
}
