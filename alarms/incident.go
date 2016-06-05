package alarms

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

// openIncident opens a new alarm incident
func (s *Service) openIncident(alarm *Alarm, incidentTypeID string, resp *http.Response, responseTime int64, errMsg string) error {
	// Begin a transaction
	tx := s.db.Begin()

	now := gorm.NowFunc()

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

		// Assign related object
		incident.Alarm = alarm

		// There should be only one open incident per alarm at any time
		err = tx.Model(new(Incident)).Where(
			"resolved_at IS NULL AND alarm_id = ? AND id != ?",
			alarm.ID,
			incident.ID,
		).UpdateColumns(Incident{
			ResolvedAt: util.TimeOrNull(&now),
			Model:      gorm.Model{UpdatedAt: now},
		}).Error
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		// Keep incidents slice up-to-date
		for i := range alarm.Incidents {
			if !alarm.Incidents[i].ResolvedAt.Valid {
				alarm.Incidents[i].ResolvedAt = util.TimeOrNull(&now)
			}
		}
		alarm.Incidents = append(alarm.Incidents, incident)

		// Send new incident push notification alert
		if alarm.PushNotificationAlerts {
			go func(a *Alarm, i *Incident) {
				endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
					uint(a.UserID.Int64),
					s.cnf.AWS.APNSPlatformApplicationARN,
				)
				if err == nil && endpoint != nil {
					_, err := s.notificationsService.PublishMessage(
						endpoint.ARN,
						fmt.Sprintf(
							newIncidentPushNotificationTemplates[i.IncidentTypeID.String],
							a.EndpointURL,
						),
						map[string]interface{}{},
					)
					if err != nil {
						logger.Errorf("Publish Message Error: %s", err.Error())
					}
				}
			}(alarm, incident)
		}

		// Send new incident notification email alert
		if alarm.EmailAlerts {
			go func(i *Incident) {
				newIncidentEmail := s.emailFactory.NewIncidentEmail(i)

				// Try to send the new incident email
				if err := s.emailService.Send(newIncidentEmail); err != nil {
					logger.Errorf("Send email error: %s", err)
					return
				}
			}(incident)
		}

		if alarm.SlackAlerts && alarm.User.SlackIncomingWebhook.Valid && alarm.User.SlackChannel.Valid {
			go func(a *Alarm) {
				// Fetch the user team
				team, _ := s.teamsService.FindTeamByMemberID(a.User.ID)

				// Get alarm limits
				alarmLimits := s.getAlarmLimits(team, a.User)

				if !alarmLimits.slackAlerts {
					return
				}

				// if err := s.slackAdapter.SendMessage(
				// 	alarm.User.SlackIncomingWebhook.String,
				// 	alarm.User.SlackChannel.String,
				// 	slackNotificationsUsername,
				// 	"TODO",
				// 	slackNotificationsEmoji,
				// ); err != nil {
				// 	logger.Errorf("Send Slack message error: %s", err)
				// 	return
				// }
			}(alarm)
		}
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
	// If the alarm state is alarmstates.OK, just return, nothing to do
	if alarm.AlarmStateID.String == alarmstates.OK {
		return nil
	}

	var err error

	// Begin a transaction
	tx := s.db.Begin()

	now := gorm.NowFunc()

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

	// Resolve open incidents
	err = tx.Model(new(Incident)).Where(
		"resolved_at IS NULL AND alarm_id = ?", alarm.ID,
	).UpdateColumns(Incident{
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

	// Send incidents resolved push notification alert
	if alarm.PushNotificationAlerts && alarmInitialState != alarmstates.InsufficientData {
		go func() {
			endpoint, err := s.notificationsService.FindEndpointByUserIDAndApplicationARN(
				alarm.User.ID,
				s.cnf.AWS.APNSPlatformApplicationARN,
			)
			if err == nil && endpoint != nil {
				_, err := s.notificationsService.PublishMessage(
					endpoint.ARN,
					fmt.Sprintf(incidentsResolvedPushNotificationTemplate, alarm.EndpointURL),
					map[string]interface{}{},
				)
				if err != nil {
					logger.Errorf("Publish Message Error: %s", err.Error())
				}
			}
		}()
	}

	// Send incidents resolved notification email alert
	if alarm.EmailAlerts && alarmInitialState != alarmstates.InsufficientData {
		go func() {
			alarmUpEmail := s.emailFactory.NewIncidentsResolvedEmail(alarm)

			// Try to send the alarm up email email
			if err := s.emailService.Send(alarmUpEmail); err != nil {
				logger.Errorf("Send email error: %s", err)
				return
			}
		}()
	}

	return nil
}

// incidentTypeCount returns aggregated count of incidents types
func (s *Service) incidentTypeCounts(user *accounts.User, alarm *Alarm, from, to *time.Time) (map[string]int, error) {
	var incitentTypeCounts = map[string]int{
		incidenttypes.Slow:    0,
		incidenttypes.Timeout: 0,
		incidenttypes.BadCode: 0,
		incidenttypes.Other:   0,
	}

	// Run aggregate count query grouped by incident type
	rows, err := s.incidentsQuery(user, alarm, nil, from, to).
		Select("incident_type_id, COUNT(*)").Group("incident_type_id").Rows()
	if err != nil {
		return incitentTypeCounts, err
	}

	// Iterate over *sql.Rows
	for rows.Next() {
		// Declare vars for copying the data from the row
		var (
			incidentTypeID string
			count          int
		)

		// Scan the data into our vars
		if err := rows.Scan(&incidentTypeID, &count); err != nil {
			return incitentTypeCounts, err
		}
		incitentTypeCounts[incidentTypeID] = count
	}

	return incitentTypeCounts, nil
}

// getUptimeDowntime returns uptime and downtime percentages
func (s *Service) getUptimeDowntime(alarm *Alarm) (float64, float64, error) {
	query := `SELECT
		COALESCE(
			EXTRACT(EPOCH FROM (SELECT NOW() - created_at FROM alarm_alarms WHERE id = ?))
				- EXTRACT(EPOCH FROM (SUM(resolved_at - created_at))),
			100
		) as uptime,
		COALESCE(
			EXTRACT(EPOCH FROM (SUM(COALESCE(resolved_at, NOW()) - created_at))),
			0
		) AS downtime
	FROM alarm_incidents WHERE alarm_id = ? AND resolved_at IS NOT NULL;`
	row := s.db.Raw(query, alarm.ID, alarm.ID).Row()

	var uptime, downtime float64
	if err := row.Scan(&uptime, &downtime); err != nil {
		return 0, 0, err
	}

	// Calculate percentages
	total := uptime + downtime
	if uptime > 0 {
		uptime = uptime / total * 100
	}
	if downtime > 0 {
		downtime = 100 - uptime
	}

	return uptime, downtime, nil
}

// incidentsCount returns a total count of incidents
func (s *Service) incidentsCount(user *accounts.User, alarm *Alarm, incidentType *string, from, to *time.Time) (int, error) {
	var count int
	err := s.incidentsQuery(user, alarm, incidentType, from, to).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// findPaginatedIncidents returns paginated incident records
// Results can optionally be filtered by user and/or alarm
func (s *Service) findPaginatedIncidents(offset, limit int, orderBy string, user *accounts.User, alarm *Alarm, incidentType *string, from, to *time.Time) ([]*Incident, error) {
	var incidents []*Incident

	// Get the pagination query
	incidentsQuery := s.incidentsQuery(user, alarm, incidentType, from, to)

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

// incidentsQuery returns a generic db query for fetching incidents
func (s *Service) incidentsQuery(user *accounts.User, alarm *Alarm, incidentType *string, from, to *time.Time) *gorm.DB {
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

	// Optionally filter by incident type
	if incidentType != nil {
		incidentsQuery = incidentsQuery.Where(Incident{
			IncidentTypeID: util.StringOrNull(*incidentType),
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
