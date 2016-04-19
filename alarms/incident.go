package alarms

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// openIncident opens a new alarm incident
func (s *Service) openIncident(alarm *Alarm, incidentTypeID string, resp *http.Response, responseTime int64, errMsg string) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Change the alarm state to alarmstates.Alarm if it isn't already
	if alarm.AlarmStateID.String != alarmstates.Alarm {
		now := gorm.NowFunc()
		err := tx.Model(alarm).UpdateColumns(Alarm{
			AlarmStateID:          util.StringOrNull(alarmstates.Alarm),
			LastDowntimeStartedAt: util.TimeOrNull(&now),
			Model: gorm.Model{UpdatedAt: now},
		}).Error
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		// Send alarm down push notification alert
		if alarm.PushNotificationAlerts {
			go func() {
				endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
					alarm.User.ID,
					s.cnf.AWS.APNSPlatformApplicationARN,
				)
				if err == nil && endpoint != nil {
					_, err := s.notificationsService.PublishMessage(
						endpoint.ARN,
						fmt.Sprintf("ALERT: %s is down", alarm.EndpointURL),
						map[string]interface{}{},
					)
					if err != nil {
						logger.Errorf("Publish Message Error: %s", err.Error())
					}
				}
			}()
		}

		// Send alarm down notification email alert
		if alarm.EmailAlerts {
			go func() {
				alarmDownEmail := s.emailFactory.NewAlarmDownEmail(alarm)

				// Try to send the alarm down email email
				if err := s.emailService.Send(alarmDownEmail); err != nil {
					logger.Errorf("Send email error: %s", err)
					return
				}
			}()
		}
	}

	// If the alarm does not have an open incident of such type yet
	if !alarm.HasOpenIncident(incidentTypeID, resp, errMsg) {
		// Fetch the incident type from the database
		incidentType, err := s.findIncidentTypeByID(incidentTypeID)
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		// Create a new incident object
		incident := NewIncident(
			alarm,
			incidentType,
			resp,
			responseTime,
			errMsg,
		)

		// Save the incident to the database
		if err := tx.Create(incident).Error; err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		alarm.Incidents = append(alarm.Incidents, incident)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// resolveIncidents resolves any open alarm incidents
func (s *Service) resolveIncidents(alarm *Alarm) error {
	var err error

	// Begin a transaction
	tx := s.db.Begin()

	now := gorm.NowFunc()

	// Change the alarm state to alarmstates.OK if it isn't already
	if alarm.AlarmStateID.String != alarmstates.OK {
		// Save the current alarm state in a variable
		alarmInitialState := alarm.AlarmStateID.String

		// Set state to alarmstates.OK and update uptime timestamp
		err = tx.Model(alarm).UpdateColumns(Alarm{
			AlarmStateID:        util.StringOrNull(alarmstates.OK),
			LastUptimeStartedAt: util.TimeOrNull(&now),
			Model:               gorm.Model{UpdatedAt: now},
		}).Error
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		// Send alarm up push notification alert
		if alarm.PushNotificationAlerts && alarmInitialState != alarmstates.InsufficientData {
			go func() {
				endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
					alarm.User.ID,
					s.cnf.AWS.APNSPlatformApplicationARN,
				)
				if err == nil && endpoint != nil {
					_, err := s.notificationsService.PublishMessage(
						endpoint.ARN,
						fmt.Sprintf("ALERT: %s is up again", alarm.EndpointURL),
						map[string]interface{}{},
					)
					if err != nil {
						logger.Errorf("Publish Message Error: %s", err.Error())
					}
				}
			}()
		}

		// Send alarm up notification email alert
		if alarm.EmailAlerts && alarmInitialState != alarmstates.InsufficientData {
			go func() {
				alarmUpEmail := s.emailFactory.NewAlarmUpEmail(alarm)

				// Try to send the alarm up email email
				if err := s.emailService.Send(alarmUpEmail); err != nil {
					logger.Errorf("Send email error: %s", err)
					return
				}
			}()
		}
	}

	// Resolve open incidents
	err = tx.Model(new(Incident)).Where("resolved_at IS NULL AND alarm_id = ?", alarm.ID).UpdateColumns(Incident{
		ResolvedAt: util.TimeOrNull(&now),
		Model:      gorm.Model{UpdatedAt: now},
	}).Error
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Make sure incidents of the passed alarm object are up-to-date
	for _, incident := range alarm.Incidents {
		incident.ResolvedAt = util.TimeOrNull(&now)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// paginatedIncidentsCount returns a total count of incidents
func (s *Service) paginatedIncidentsCount(user *accounts.User, alarm *Alarm, from, to *time.Time) (int, error) {
	var count int
	err := s.paginatedIncidentsQuery(user, alarm, from, to).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedIncidents returns paginated incident records
// Results can optionally be filtered by user and/or alarm
func (s *Service) findPaginatedIncidents(offset, limit int, orderBy string, user *accounts.User, alarm *Alarm, from, to *time.Time) ([]*Incident, error) {
	var incidents []*Incident

	// Get the pagination query
	incidentsQuery := s.paginatedIncidentsQuery(user, alarm, from, to)

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
func (s *Service) paginatedIncidentsQuery(user *accounts.User, alarm *Alarm, from, to *time.Time) *gorm.DB {
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

	// Optionally filter incidents older than from
	if from != nil {
		incidentsQuery = incidentsQuery.Where("created_at >= ?", from)
	}

	// Optionally filter incidents younger than to
	if to != nil {
		incidentsQuery = incidentsQuery.Where("created_at <= ?", to)
	}

	return incidentsQuery
}
