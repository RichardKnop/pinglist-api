package alarms

import (
	"errors"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrAlarmNotFound ...
	ErrAlarmNotFound = errors.New("Alarm not found")
	// ErrMaxAlarmsLimitReached ...
	ErrMaxAlarmsLimitReached = errors.New("Max alarms limit reached")
)

// HasOpenIncident returns true if the alarm already has such open incident
func (a *Alarm) HasOpenIncident(theType string, resp *http.Response, errMsg string) bool {
	for _, incident := range a.Incidents {
		// If incident is resolved, continue the loop
		if incident.ResolvedAt.Valid {
			continue
		}

		// If incident is of a different type, continue the loop
		if incident.IncidentTypeID.String != theType {
			continue
		}

		isBadCode := incident.IncidentTypeID.String == incidenttypes.BadCode

		// For other than bad code incidents, we compare the error message
		if !isBadCode {
			if incident.ErrorMessage.String == errMsg {
				return true
			}
		}

		// For bad code incidents, we compare the status code
		if isBadCode {
			if resp != nil && incident.HTTPCode.Valid && int64(resp.StatusCode) == incident.HTTPCode.Int64 {
				return true
			}
		}
	}

	return false
}

// FindAlarmByID looks up an alarm by ID and returns it
func (s *Service) FindAlarmByID(alarmID uint) (*Alarm, error) {
	// Fetch the alarm from the database
	alarm := new(Alarm)
	notFound := s.db.Preload("User.OauthUser").Preload("Incidents", "resolved_at IS NULL").
		First(alarm, alarmID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAlarmNotFound
	}

	return alarm, nil
}

// createAlarm creates a new alarm
func (s *Service) createAlarm(user *accounts.User, alarmRequest *AlarmRequest) (*Alarm, error) {
	if alarmRequest.Active {
		// Limit maximum number of active alarms per user
		alarmsCount, err := s.userActiveAlarmsCount(user)
		if err != nil {
			return nil, err
		}
		if alarmsCount+1 > s.getUserMaxAlarms(user) {
			return nil, ErrMaxAlarmsLimitReached
		}
	}

	// Fetch the region from the database
	region, err := s.findRegionByID(alarmRequest.Region)
	if err != nil {
		return nil, err
	}

	// Fetch the initial alarm state from the database
	alarmState, err := s.findAlarmStateByID(alarmstates.InsufficientData)
	if err != nil {
		return nil, err
	}

	// Create a new alarm object
	alarm := NewAlarm(user, region, alarmState, alarmRequest)

	// Save the alarm to the database
	if err := s.db.Create(alarm).Error; err != nil {
		return nil, err
	}

	return alarm, nil
}

// updateAlarm updates an existing alarm
func (s *Service) updateAlarm(alarm *Alarm, alarmRequest *AlarmRequest) error {
	if !alarm.Active && alarmRequest.Active {
		// Limit maximum number of active alarms per user
		alarmsCount, err := s.userActiveAlarmsCount(alarm.User)
		if err != nil {
			return err
		}
		if alarmsCount+1 > s.getUserMaxAlarms(alarm.User) {
			return ErrMaxAlarmsLimitReached
		}
	}

	// Fetch the region from the database
	region, err := s.findRegionByID(alarmRequest.Region)
	if err != nil {
		return err
	}

	// Update the alarm (need to use map here because active field might be
	// changing to false which would not work with struct)
	if err := s.db.Model(alarm).UpdateColumns(map[string]interface{}{
		"region_id":                region.ID,
		"endpoint_url":             alarmRequest.EndpointURL,
		"expected_http_code":       alarmRequest.ExpectedHTTPCode,
		"interval":                 alarmRequest.Interval,
		"email_alerts":             alarmRequest.EmailAlerts,
		"push_notification_alerts": alarmRequest.PushNotificationAlerts,
		"active":                   alarmRequest.Active,
		"updated_at":               time.Now(),
	}).Error; err != nil {
		return err
	}

	// Make sure the alarm region is up-to-date
	alarm.Region = region

	return nil
}

// userActiveAlarmsCount counts active alarms of a user
func (s *Service) userActiveAlarmsCount(user *accounts.User) (int, error) {
	var count int
	err := s.db.Model(new(Alarm)).Where("user_id = ?", user.ID).
		Where("active = ?", true).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
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
