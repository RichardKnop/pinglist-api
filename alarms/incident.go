package alarms

import (
	"net/http"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// openIncident opens a new alarm incident
func (s *Service) openIncident(alarm *Alarm, theType string, resp *http.Response) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Change the alarm state to alarmstates.Alarm if it isn't already
	if alarm.State != alarmstates.Alarm {
		err := tx.Model(alarm).UpdateColumn("state", alarmstates.Alarm).Error
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
	}

	var incident *Incident

	// If the alarm does not have an open incident of such type yet
	if !alarm.HasOpenIncident(theType, resp) {
		// Create a new incident object
		incident = newIncident(
			alarm,
			theType,
			resp,
		)

		// Save the incident to the database
		if err := tx.Create(incident).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Make sure to keep the passed alarm object up-to-date
	// if a new incident was opened
	if incident != nil {
		alarm.Incidents = append(alarm.Incidents, incident)
	}

	return nil
}

// resolveIncidentsTx resolves any open alarm incidents inside a transaction
func (s *Service) resolveIncidentsTx(db *gorm.DB, alarm *Alarm) error {
	var err error

	// Change alarm state to alarmstates.OK
	err = db.Model(alarm).UpdateColumn("state", alarmstates.OK).Error
	if err != nil {
		return err
	}

	// Resolve incidents
	now := gorm.NowFunc()
	err = db.Model(new(Incident)).Where(Incident{
		AlarmID: util.PositiveIntOrNull(int64(alarm.ID)),
	}).UpdateColumn("resolved_at", util.TimeOrNull(&now)).Error
	if err != nil {
		return err
	}

	// Make sure incidents of the passed alarm object are up-to-date
	for _, incident := range alarm.Incidents {
		incident.ResolvedAt = util.TimeOrNull(&now)
	}

	return nil
}

// paginatedIncidentsCount returns a total count of incidents
func (s *Service) paginatedIncidentsCount(user *accounts.User, alarm *Alarm) (int, error) {
	var count int
	if err := s.paginatedIncidentsQuery(user, alarm).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedIncidents returns paginated incident records
// Results can optionally be filtered by user and/or alarm
func (s *Service) findPaginatedIncidents(offset, limit int, orderBy string, user *accounts.User, alarm *Alarm) ([]*Incident, error) {
	var incidents []*Incident

	// Get the pagination query
	incidentsQuery := s.paginatedIncidentsQuery(user, alarm)

	// Default ordering
	if orderBy == "" {
		orderBy = "id"
	}

	// Retrieve paginated results from the database
	err := incidentsQuery.Offset(offset).Limit(limit).Order(orderBy).
		Preload("Alarm.User").Find(&incidents).Error
	if err != nil {
		return incidents, err
	}

	return incidents, nil
}

// paginatedIncidentsQuery returns a db query for paginated incidents
func (s *Service) paginatedIncidentsQuery(user *accounts.User, alarm *Alarm) *gorm.DB {
	// Basic query
	incidentsQuery := s.db.Model(new(Incident))

	// Optionally filter by user
	if user != nil {
		incidentsQuery = incidentsQuery.
			Joins("inner join alarm_alarms on alarm_alarms.id = alarm_incidents.alarm_id").
			Where("alarm_alarms.user_id = ?", user.ID)
	}

	// Optionally filter by alarm
	if alarm != nil {
		incidentsQuery = incidentsQuery.Where(Incident{
			AlarmID: util.PositiveIntOrNull(int64(alarm.ID)),
		})
	}

	return incidentsQuery
}
