package alarms

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrCheckAlreadyTriggered ...
	ErrCheckAlreadyTriggered = errors.New("Alarm check has already been trigerred")
)

// GetAlarmsToCheck returns alarms that should be checked
func (s *Service) GetAlarmsToCheck(now time.Time) ([]*Alarm, error) {
	var alarms []*Alarm

	watermarkCondition := "watermark IS NULL OR watermark + interval '1 second' * interval > ?"
	err := s.db.Where("active = ?", true).Where(watermarkCondition, now).Order("id").Find(&alarms).Error
	if err != nil {
		return alarms, err
	}

	return alarms, nil
}

// CheckAlarm performs an alarm check
func (s *Service) CheckAlarm(alarmID uint, watermark time.Time) error {
	// Fetch the alarm
	alarm := new(Alarm)
	if s.db.First(alarm, alarmID).RecordNotFound() {
		return ErrAlarmNotFound
	}

	// Idempotency check
	if alarm.Watermark.Time.After(watermark) {
		return ErrCheckAlreadyTriggered
	}

	// Prepare a request
	req, err := http.NewRequest("GET", alarm.EndpointURL, nil)
	if err != nil {
		return err
	}

	// Update the watermark
	newWatermark := gorm.NowFunc()
	err = s.db.Model(alarm).UpdateColumn("watermark", util.TimeOrNull(&newWatermark)).Error
	if err != nil {
		return err
	}

	// Make the request
	start := time.Now()
	resp, err := s.client.Do(req)
	elapsed := time.Since(start)

	// The request timed out
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return s.openIncident(alarm, incidenttypes.Timeout, resp)
	}

	// The request failed due to any other error
	if err != nil {
		return s.openIncident(alarm, incidenttypes.Other, resp)
	}

	defer resp.Body.Close()

	// The request returned a response with a bad status code
	if resp.StatusCode != int(alarm.ExpectedHTTPCode) {
		return s.openIncident(alarm, incidenttypes.BadCode, resp)
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Resolve any open incidents
	if err := s.resolveIncidentsTx(tx, alarm); err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Create a new result object
	result := newResult(
		getSubtableName(ResultParentTableName, start),
		alarm,
		newWatermark,
		elapsed.Nanoseconds(),
	)

	// Save the result to the database
	if err := tx.Create(result).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Make sure to keep the passed alarm object up-to-date
	alarm.Results = append(alarm.Results, result)

	return nil
}
