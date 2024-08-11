package services

import (
	"sync"
	"time"
)

type Location struct {
	Name       string
	Latitude   string
	Longtitude string
	Timestamp  time.Time
}

type LocationService struct {
	locations map[string]Location
	mu        sync.RWMutex
}

func NewLocationService() *LocationService {
	return &LocationService{
		locations: make(map[string]Location),
	}
}

func (s *LocationService) SetLocation(openID, name, lat, lon string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.locations[openID] = Location{
		Name:       name,
		Latitude:   lat,
		Longtitude: lon,
		Timestamp:  time.Now(),
	}
}

func (s *LocationService) GetLocation(openID string) (Location, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	location, ok := s.locations[openID]
	return location, ok
}

func (s *LocationService) SaveLocation(openID, name, lat, lon string) {
	s.SetLocation(openID, name, lat, lon)
}
