package alarms

import (
	"errors"
)

var (
	// ErrAlarmStateNotFound ...
	ErrAlarmStateNotFound = errors.New("Alarm state not found")
)

// findAlarmStateByID looks up an alarm state by ID and returns it
func (s *Service) findAlarmStateByID(id string) (*AlarmState, error) {
	alarmState := new(AlarmState)
	if s.db.Where("id = ?", id).First(alarmState).RecordNotFound() {
		return nil, ErrAlarmStateNotFound
	}
	return alarmState, nil
}
