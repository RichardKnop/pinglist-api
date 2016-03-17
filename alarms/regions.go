package alarms

import (
	"errors"
)

var (
	// ErrRegionNotFound ...
	ErrRegionNotFound = errors.New("Region not found")
)

// findRegionByID looks up a region by ID and returns it
func (s *Service) findRegionByID(id string) (*Region, error) {
	region := new(Region)
	if s.db.Where("id = ?", id).First(region).RecordNotFound() {
		return nil, ErrRegionNotFound
	}
	return region, nil
}
