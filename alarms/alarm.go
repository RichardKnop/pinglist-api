package alarms

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
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
func (a *Alarm) HasOpenIncident(theType string, resp *http.Response) bool {
	for _, incident := range a.Incidents {
		if incident.ResolvedAt.Valid || incident.IncidentTypeID.String != theType {
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
		return nil, ErrAlarmNotFound
	}

	return alarm, nil
}

// createAlarm creates a new alarm
func (s *Service) createAlarm(user *accounts.User, alarmRequest *AlarmRequest) (*Alarm, error) {
	var maxAlarms int

	// If user is in a free trial, allow one alarm
	if user.IsInFreeTrial() {
		maxAlarms = 1
	}

	// Fetch active user subscription
	subscription, err := s.subscriptionsService.FindActiveUserSubscription(user.ID)

	// If subscribed, take max allowed alarms from the subscription plan
	if err == nil && subscription != nil {
		maxAlarms = int(subscription.Plan.MaxAlarms)
	}

	// Count how many alarms user already has
	var countAlarms int
	s.db.Model(new(Alarm)).Where("user_id = ?", user.ID).Count(&countAlarms)

	// Limit alarms to max number defined above
	if countAlarms+1 > maxAlarms {
		return nil, ErrMaxAlarmsLimitReached
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
	alarm := newAlarm(user, region, alarmState, alarmRequest)

	// Save the alarm to the database
	if err := s.db.Create(alarm).Error; err != nil {
		return nil, err
	}

	return alarm, nil
}

// updateAlarm updates an existing alarm
func (s *Service) updateAlarm(alarm *Alarm, alarmRequest *AlarmRequest) error {
	// Fetch the region from the database
	region, err := s.findRegionByID(alarmRequest.Region)
	if err != nil {
		return err
	}

	// Update the alarm
	if err := s.db.Model(alarm).UpdateColumns(Alarm{
		RegionID:         util.StringOrNull(region.ID),
		EndpointURL:      alarmRequest.EndpointURL,
		ExpectedHTTPCode: alarmRequest.ExpectedHTTPCode,
		Interval:         alarmRequest.Interval,
		Active:           alarmRequest.Active,
	}).Error; err != nil {
		return err
	}

	// Make sure the alarm region is up-to-date
	alarm.Region = region

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
