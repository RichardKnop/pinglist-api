package alarms

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestGetAlarmsToCheck() {
	var (
		alarms    []*Alarm
		err       error
		watermark time.Time
		interval  = uint(60)
	)

	// Deactivate all alarms
	err = suite.service.db.Model(new(Alarm)).UpdateColumn("active", false).Error
	assert.NoError(suite.T(), err, "Deactivating alarms failed")

	// First, let's try with no active alarms
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarms
	assert.Equal(suite.T(), 0, len(alarms))

	// Now insert an active test alarm with watermark + interval <= now
	watermark = time.Now().Add(-time.Duration(interval+1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		EndpointURL:      "http://foo",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		Interval:         interval,
		Active:           true,
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Try again
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 0 alarms
	assert.Equal(suite.T(), 0, len(alarms))

	// Now insert an active test alarm with watermark + interval > now
	watermark = time.Now().Add(-time.Duration(interval-1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		EndpointURL:      "http://bar",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		Interval:         interval,
		Active:           true,
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Try again
	alarms, err = suite.service.GetAlarmsToCheck(time.Now())

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 alarms
	assert.Equal(suite.T(), 1, len(alarms))
	assert.Equal(suite.T(), "http://bar", alarms[0].EndpointURL)
}

func (suite *AlarmsTestSuite) TestAlarmCheck() {
	var (
		alarm  = suite.alarms[2]
		err    error
		server *httptest.Server
		client *http.Client
		start  time.Time
	)

	// Partition the results table
	err = suite.service.PartitionTable(ResultParentTableName, time.Now())
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// First, let's test a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.State)

	// 0 incidents, 1 result
	assert.Equal(suite.T(), 0, len(alarm.Incidents))
	assert.Equal(suite.T(), 1, len(alarm.Results))

	// New result
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Results[0].AlarmID.Int64))
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Results[0].Timestamp.Format("2006-01-02T15:04:05Z"),
	)
	assert.True(suite.T(), alarm.Results[0].RequestTime > 0)

	// Second, let's test a timeout
	server, client = testServerTimeout()
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.State)

	// 1 incident, 1 result
	assert.Equal(suite.T(), 1, len(alarm.Incidents))
	assert.Equal(suite.T(), 1, len(alarm.Results))

	// New incident
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Incidents[0].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[0].Type)
	assert.False(suite.T(), alarm.Incidents[0].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Third, let's test a bad code
	server, client = testServer(&http.Response{StatusCode: 500})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(suite.alarms[2].ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.State)

	// 2 incidents, 1 result
	assert.Equal(suite.T(), 2, len(alarm.Incidents))
	assert.Equal(suite.T(), 1, len(alarm.Results))

	// New incident
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Incidents[1].AlarmID.Int64))
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[1].Type)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[1].HTTPCode.Int64)
	assert.Equal(suite.T(), "", alarm.Incidents[1].Response.String)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Finally, let's test a return to a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(suite.alarms[2].ID, alarm.Watermark.Time)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.State)

	// 2 incidents, 2 results
	assert.Equal(suite.T(), 2, len(alarm.Incidents))
	assert.Equal(suite.T(), 2, len(alarm.Results))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}

	// New result
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Results[1].AlarmID.Int64))
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Results[1].Timestamp.Format("2006-01-02T15:04:05Z"),
	)
	assert.True(suite.T(), alarm.Results[1].RequestTime > 0)
}

func (suite *AlarmsTestSuite) TestAlarmCheckIdempotency() {
	// Fetch the alarm
	alarm := new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Partition the results table
	err := suite.service.PartitionTable(ResultParentTableName, time.Now())
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Prepare test server and client
	server, client := testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start := time.Now()
	gorm.NowFunc = func() time.Time {
		return start
	}

	concurrency := 4

	// Trigger multiple parallel alarm checks with the same watermark
	errChan := make(chan error)
	for i := 0; i < concurrency; i++ {
		go suite.alarmCheckWrapper(alarm.ID, alarm.Watermark.Time, errChan)
		time.Sleep(10 * time.Millisecond)
	}

	// Receive start times and errors from goroutines
	var (
		successful int
		errs       []error
	)
	for i := 0; i < concurrency; i++ {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		} else {
			successful++
		}
	}

	// One alarm check should have gone through
	assert.Equal(suite.T(), 1, successful)

	// 3 alarm checks should have been stopped by the idempotency check
	assert.Equal(suite.T(), 3, len(errs))
	for _, err := range errs {
		assert.Equal(suite.T(), errCheckAlreadyTriggered, err)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		Preload("Results").First(alarm, suite.alarms[2].ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.State)

	// 0 incidents, 1 result
	assert.Equal(suite.T(), 0, len(alarm.Incidents))
	assert.Equal(suite.T(), 1, len(alarm.Results))

	// New result
	assert.Equal(suite.T(), suite.alarms[2].ID, uint(alarm.Results[0].AlarmID.Int64))
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Results[0].Timestamp.Format("2006-01-02T15:04:05Z"),
	)
	assert.True(suite.T(), alarm.Results[0].RequestTime > 0)
}

func (suite *AlarmsTestSuite) alarmCheckWrapper(alarmID uint, watermark time.Time, errChan chan error) {
	errChan <- suite.service.CheckAlarm(alarmID, watermark)
}
