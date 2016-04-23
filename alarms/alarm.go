package alarms

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// MinIntervalLimit limits alarm check interval to a sensible smallest period
	MinIntervalLimit = uint(60)
	// MaxResponseTimeLimit limits max response time to a sensible biggest value
	MaxResponseTimeLimit = uint(10000)

	// ErrAlarmNotFound ...
	ErrAlarmNotFound = errors.New("Alarm not found")
	// ErrMaxAlarmsLimitReached ...
	ErrMaxAlarmsLimitReached = errors.New("Max alarms limit reached")
	// ErrIntervalTooSmall ...
	ErrIntervalTooSmall = fmt.Errorf("Minimal interval is %d seconds", MinIntervalLimit)
	// ErrMaxResponseTimeTooBig ...
	ErrMaxResponseTimeTooBig = fmt.Errorf("Max response time cannot be greater than %d ms", MaxResponseTimeLimit)
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
	// Limit active alarms to the max number defined as per subscription plan
	if alarmRequest.Active {
		// Fetch the user team
		team, _ := s.teamsService.FindTeamByMemberID(user.ID)

		// Count alarms and calculate max limit
		var (
			alarmsCount = s.countActiveAlarms(team, user)
			maxAlarms   = s.getMaxAlarms(team, user)
		)

		// Check the alarm limit
		if alarmsCount+1 > maxAlarms {
			return nil, ErrMaxAlarmsLimitReached
		}
	}

	// Limit interval to a sensible smallest period (MinInterval constant)
	if alarmRequest.Interval < MinIntervalLimit {
		return nil, ErrIntervalTooSmall
	}

	// Limit max response time to a sensible biggest value (ErrTimeoutTooBig constact)
	if alarmRequest.MaxResponseTime > MaxResponseTimeLimit {
		return nil, ErrMaxResponseTimeTooBig
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
	// Limit active alarms to the max number defined as per subscription plan
	if !alarm.Active && alarmRequest.Active {
		// Fetch the user team
		team, _ := s.teamsService.FindTeamByMemberID(alarm.User.ID)

		// Count alarms and calculate max limit
		var (
			alarmsCount = s.countActiveAlarms(team, alarm.User)
			maxAlarms   = s.getMaxAlarms(team, alarm.User)
		)

		// Check the alarm limit
		if alarmsCount+1 > maxAlarms {
			return ErrMaxAlarmsLimitReached
		}
	}

	// Limit interval to a sensible smallest period (MinInterval constant)
	if alarmRequest.Interval < MinIntervalLimit {
		return ErrIntervalTooSmall
	}

	// Limit max response time to a sensible biggest value (ErrTimeoutTooBig constact)
	if alarmRequest.MaxResponseTime > MaxResponseTimeLimit {
		return ErrMaxResponseTimeTooBig
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
		"max_response_time":        alarmRequest.MaxResponseTime,
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

// alarmsCount returns a total count of alarms
// Can be optionally filtered by user
func (s *Service) alarmsCount(user *accounts.User) (int, error) {
	var count int
	if err := s.alarmsQuery(user).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedAlarms returns paginated alarm records
// Results can optionally be filtered by user
func (s *Service) findPaginatedAlarms(offset, limit int, orderBy string, user *accounts.User) ([]*Alarm, error) {
	var alarms []*Alarm

	// Get the pagination query
	alarmsQuery := s.alarmsQuery(user)

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

// alarmsQuery returns a generic db query for fetching alarms
func (s *Service) alarmsQuery(user *accounts.User) *gorm.DB {
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
