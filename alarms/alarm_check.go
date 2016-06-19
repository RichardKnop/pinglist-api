package alarms

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrCheckAlreadyTriggered ...
	ErrCheckAlreadyTriggered = errors.New("Alarm check has already been trigerred")

	// AlarmCheckTimeout defines how long to wait before considering alarm check timed out
	AlarmCheckTimeout = 10 * time.Second
)

// GetAlarmsToCheck returns IDs of alarms that should be checked
func (s *Service) GetAlarmsToCheck(now time.Time) ([]uint, error) {
	var alarmIDs []uint
	query := `SELECT * FROM (
		SELECT
		a.id,
		COALESCE(GREATEST(p.max_alarms, p2.max_alarms), 1) max_alarms,
		DENSE_RANK() OVER (PARTITION BY COALESCE(CAST(s.id AS TEXT), ou.username) ORDER BY a.id ASC) AS rank
		FROM alarm_alarms a
		INNER JOIN account_users u ON u.id = a.user_id
		INNER JOIN oauth_users ou ON ou.id = u.oauth_user_id
		LEFT JOIN subscription_customers c ON c.user_id = a.user_id
		LEFT JOIN subscription_subscriptions s ON s.customer_id = c.id AND s.period_end > ?
		LEFT JOIN subscription_plans p ON p.id = s.plan_id
		LEFT JOIN team_team_members tm ON tm.user_id = u.id
		LEFT JOIN team_teams t ON t.id = tm.team_id
		LEFT JOIN subscription_customers c2 ON c2.user_id = t.owner_id
		LEFT JOIN subscription_subscriptions s2 ON s2.customer_id = c2.id AND s2.period_end > ?
		LEFT JOIN subscription_plans p2 ON p2.id = s2.plan_id
		WHERE
		(watermark IS NULL OR watermark + interval '1 second' * a.interval < ?) AND active=true
		ORDER BY s.id, a.user_id, rank
	) t WHERE rank <= max_alarms;`
	rows, err := s.db.Raw(query, now, now, now).Rows() // (*sql.Rows, error)
	if err != nil {
		return alarmIDs, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			alarmID   uint
			maxAlarms uint
			rank      uint
		)
		if err := rows.Scan(&alarmID, &maxAlarms, &rank); err != nil {
			return alarmIDs, err
		}
		alarmIDs = append(alarmIDs, alarmID)
	}
	return alarmIDs, nil
}

// CheckAlarm performs an alarm check
func (s *Service) CheckAlarm(alarmID uint, watermark time.Time) error {
	// Fetch the alarm
	alarm, err := s.FindAlarmByID(alarmID)
	if err != nil {
		return err
	}

	// Idempotency check
	if alarm.Watermark.Time.After(watermark) {
		return ErrCheckAlreadyTriggered
	}

	// Start with the default request URL
	var requestURL = alarm.EndpointURL

	// If we are going to proxy the request through a remote server
	if alarm.Region.ID != s.cnf.AWS.Region {
		// Parse the proxy URL
		proxyURL, err := url.ParseRequestURI(alarm.Region.ProxyURL.String)
		if err != nil {
			return err
		}

		// Create a proxy server
		frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(proxyURL))
		defer frontendProxy.Close()

		// Set the request URL to use the proxy server
		requestURL = fmt.Sprintf(
			"%s?request_url=%s",
			frontendProxy.URL,
			alarm.EndpointURL,
		)
	}

	// Prepare a request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return err
	}

	// Update the watermark
	newWatermark := gorm.NowFunc()
	err = s.db.Model(alarm).UpdateColumns(Alarm{
		Watermark: util.TimeOrNull(&newWatermark),
		Model:     gorm.Model{UpdatedAt: newWatermark},
	}).Error
	if err != nil {
		return err
	}

	// Make the request
	start := gorm.NowFunc()
	resp, err := s.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	elapsed := time.Since(start)

	var (
		incidentType string
		errMsg       string
	)
	if e, ok := err.(net.Error); ok && e.Timeout() {
		// The response timed out
		incidentType = incidenttypes.Timeout
		errMsg = err.Error()
	} else if err != nil {
		// The request failed due to any other error
		incidentType = incidenttypes.Other
		errMsg = err.Error()
	} else if resp.StatusCode != int(alarm.ExpectedHTTPCode) {
		// The request returned a response with a bad status code
		incidentType = incidenttypes.BadCode
	} else if uint(elapsed.Nanoseconds()/1000000) > alarm.MaxResponseTime {
		// The response was too slow
		incidentType = incidenttypes.Slow
	}

	if incidentType != "" {
		// Open a new incident
		if err := s.openIncident(
			alarm,
			incidentType,
			resp,
			elapsed.Nanoseconds(),
			errMsg,
		); err != nil {
			return err
		}
	} else {
		// Resolve any open incidents
		if err := s.resolveIncidents(alarm); err != nil {
			return err
		}
	}

	// Log the response time metric
	return s.metricsService.LogResponseTime(start, alarm.ID, elapsed.Nanoseconds())
}
