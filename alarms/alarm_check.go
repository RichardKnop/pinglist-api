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

// GetAlarmsToCheck returns alarms that should be checked
func (s *Service) GetAlarmsToCheck(now time.Time) ([]*Alarm, error) {
	var alarms []*Alarm

	watermarkCondition := "watermark IS NULL OR watermark + interval '1 second' * interval < ?"
	err := s.db.Where("active = ?", true).Where(watermarkCondition, now).
		Order("id").Find(&alarms).Error
	if err != nil {
		return alarms, err
	}

	return alarms, nil
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
	elapsed := time.Since(start)

	// The response timed out
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return s.openIncident(
			alarm,
			incidenttypes.Timeout,
			nil, // response
			0,   // response time
			err.Error(),
		)
	}

	// The request failed due to any other error
	if err != nil {
		return s.openIncident(
			alarm,
			incidenttypes.Other,
			nil, // response
			0,   // response time
			err.Error(),
		)
	}

	defer resp.Body.Close()

	// The request returned a response with a bad status code
	if resp.StatusCode != int(alarm.ExpectedHTTPCode) {
		return s.openIncident(
			alarm,
			incidenttypes.BadCode,
			resp,
			elapsed.Nanoseconds(),
			"",
		)
	}

	// The response was too slow
	if uint(elapsed.Nanoseconds()/1000000) > alarm.MaxResponseTime {
		return s.openIncident(
			alarm,
			incidenttypes.SlowResponse,
			resp,
			elapsed.Nanoseconds(),
			"",
		)
	}

	// Resolve any open incidents
	if err := s.resolveIncidents(alarm); err != nil {
		return err
	}

	// Log the response time metric
	err = s.metricsService.LogResponseTime(start, alarm.ID, elapsed.Nanoseconds())
	if err != nil {
		return err
	}

	return nil
}
