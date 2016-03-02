package alarms

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	errAlarmNotFound = errors.New("Alarm not found")
)

// HasOpenIncident returns true if the alarm already has such open incident
func (a *Alarm) HasOpenIncident(theType string, resp *http.Response) bool {
	for _, incident := range a.Incidents {
		if incident.ResolvedAt.Valid || incident.Type != theType {
			continue
		}

		if resp == nil {
			return true
		}

		if incident.HTTPCode.Valid && int64(resp.StatusCode) == incident.HTTPCode.Int64 {
			return true
		}
	}

	return false
}

// FindAlarmByID looks up an alarm by ID and returns it
func (s *Service) FindAlarmByID(alarmID uint) (*Alarm, error) {
	// Fetch the alarm from the database
	alarm := new(Alarm)
	notFound := s.db.Preload("User").Preload("Incidents", "resolved_at IS NULL").
		First(alarm, alarmID).RecordNotFound()

	// Not found
	if notFound {
		return nil, errAlarmNotFound
	}

	return alarm, nil
}

// createAlarm creates a new alarm
func (s *Service) createAlarm(user *accounts.User, alarmRequest *AlarmRequest) (*Alarm, error) {
	// Create a new alarm object
	alarm := newAlarm(user, alarmRequest)

	// Save the alarm to the database
	if err := s.db.Create(alarm).Error; err != nil {
		return nil, err
	}

	return alarm, nil
}

// updateAlarm updates an existing alarm
func (s *Service) updateAlarm(alarm *Alarm, alarmRequest *AlarmRequest) error {
	// Update the alarm
	if err := s.db.Model(alarm).UpdateColumns(Alarm{
		EndpointURL:      alarmRequest.EndpointURL,
		ExpectedHTTPCode: alarmRequest.ExpectedHTTPCode,
		Interval:         alarmRequest.Interval,
		Active:           alarmRequest.Active,
	}).Error; err != nil {
		return err
	}

	return nil
}

// paginatedAlarmsCount returns a total count of alarms
// Can be optionally filtered by user
func (s *Service) paginatedAlarmsCount(user *accounts.User) (int, error) {
	var count int
	if err := s.paginatedAlarmsQuery(user).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedAlarms returns paginated alarm records
// Results can optionally be filtered by user
func (s *Service) findPaginatedAlarms(offset, limit int, orderBy string, user *accounts.User) ([]*Alarm, error) {
	var alarms []*Alarm

	// Get the pagination query
	alarmsQuery := s.paginatedAlarmsQuery(user)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := alarmsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Preload("User").Preload("Incidents").Find(&alarms).Error
	if err != nil {
		return alarms, err
	}

	return alarms, nil
}

// paginatedAlarmsQuery returns a db query for paginated alarms
func (s *Service) paginatedAlarmsQuery(user *accounts.User) *gorm.DB {
	// Basic query
	alarmsQuery := s.db.Model(new(Alarm))

	// Optionally filter by user
	if user != nil {
		alarmsQuery = alarmsQuery.Where(Alarm{
			UserID: util.PositiveIntOrNull(int64(user.ID)),
		})
	}

	return alarmsQuery
}
