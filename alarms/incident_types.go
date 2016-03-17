package alarms

import (
	"errors"
)

var (
	// ErrIncidentTypeNotFound ...
	ErrIncidentTypeNotFound = errors.New("Incident type not found")
)

// findIncidentTypeByID looks up an incident type by ID and returns it
func (s *Service) findIncidentTypeByID(id string) (*IncidentType, error) {
	incidentType := new(IncidentType)
	if s.db.Where("id = ?", id).First(incidentType).RecordNotFound() {
		return nil, ErrIncidentTypeNotFound
	}
	return incidentType, nil
}
