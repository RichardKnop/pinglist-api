package alarms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms/alarmstates"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/alarms/regions"
	"github.com/RichardKnop/pinglist-api/notifications"
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

	// Now insert an active test alarm with watermark + interval >= now
	watermark = time.Now().Add(-time.Duration(interval-1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://foo",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
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

	// Now insert an active test alarm with watermark + interval < now
	watermark = time.Now().Add(-time.Duration(interval+1) * time.Second)
	err = suite.db.Create(&Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://bar",
		Watermark:        util.TimeOrNull(&watermark),
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
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
		testAlarm, alarm *Alarm
		err              error
		server           *httptest.Server
		client           *http.Client
		start            time.Time
	)

	// Insert a test alarm
	testAlarm = &Alarm{
		User:                   suite.users[1],
		Region:                 &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:             &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:            "http://foobar",
		ExpectedHTTPCode:       200,
		MaxResponseTime:        1000,
		Interval:               60,
		EmailAlerts:            true,
		PushNotificationAlerts: true,
		Active:                 true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Fetch the alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// First, let's test a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Second, let's test a slow response
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	assert.NoError(
		suite.T(),
		suite.db.Model(alarm).UpdateColumn("max_response_time", 0).Error,
		"Updating max_response_time to 0 failed",
	)
	suite.mockAlarmDownEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s is down", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)
	assert.NoError(
		suite.T(),
		suite.db.Model(alarm).UpdateColumn("max_response_time", 1000).Error,
		"Updating max_response_time back to 1000 failed",
	)

	// Sleep for the email and push notification goroutines to finish
	time.Sleep(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status changed to Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 1 incident
	assert.Equal(suite.T(), 1, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.SlowResponse, alarm.Incidents[0].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(200), alarm.Incidents[0].HTTPCode.Int64)
	assert.True(suite.T(), alarm.Incidents[0].ResponseTime.Valid)
	assert.True(suite.T(), alarm.Incidents[0].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[0].ResolvedAt.Valid)

	// Third, let's test a timeout
	server, client = testServerTimeout()
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 2 incidents
	assert.Equal(suite.T(), 2, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.Timeout, alarm.Incidents[1].IncidentTypeID.String)
	assert.False(suite.T(), alarm.Incidents[1].HTTPCode.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResponseTime.Valid)
	assert.False(suite.T(), alarm.Incidents[1].Response.Valid)
	assert.True(suite.T(), alarm.Incidents[1].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[1].ResolvedAt.Valid)

	// Next, let's test a bad code
	server, client = testServer(&http.Response{StatusCode: 500})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status still Alarm
	assert.Equal(suite.T(), alarmstates.Alarm, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// New incident
	assert.Equal(suite.T(), incidenttypes.BadCode, alarm.Incidents[2].IncidentTypeID.String)
	assert.Equal(suite.T(), int64(500), alarm.Incidents[2].HTTPCode.Int64)
	assert.True(suite.T(), alarm.Incidents[2].ResponseTime.Valid)
	assert.True(suite.T(), alarm.Incidents[2].Response.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ErrorMessage.Valid)
	assert.False(suite.T(), alarm.Incidents[2].ResolvedAt.Valid)

	// Finally, let's test a return to a successful alarm check
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}
	suite.mockAlarmUpEmail()
	suite.mockFindEndpointByUserIDAndApplicationARN(
		alarm.User.ID,
		suite.service.cnf.AWS.APNSPlatformApplicationARN,
		&notifications.Endpoint{ARN: "endpoint_arn"},
		nil,
	)
	suite.mockPublishMessage(
		"endpoint_arn",
		fmt.Sprintf("ALERT: %s is up again", alarm.EndpointURL),
		map[string]interface{}{},
		"message_id",
		nil,
	)
	suite.mockLogResponseTime(start, alarm.ID, nil)
	err = suite.service.CheckAlarm(alarm.ID, alarm.Watermark.Time)

	// Sleep for the email & push notification goroutines to finish
	time.Sleep(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status back to OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 3 incidents
	assert.Equal(suite.T(), 3, len(alarm.Incidents))

	// Resolved incidents
	for _, incident := range alarm.Incidents {
		assert.True(suite.T(), incident.ResolvedAt.Valid)
	}
}

func (suite *AlarmsTestSuite) TestAlarmCheckIdempotency() {
	var (
		testAlarm, alarm *Alarm
		err              error
		server           *httptest.Server
		client           *http.Client
		start            time.Time
	)

	// Insert a test alarm
	testAlarm = &Alarm{
		User:             suite.users[1],
		Region:           &Region{ID: regions.USWest2, Name: "US West (Oregon)"},
		AlarmState:       &AlarmState{ID: alarmstates.InsufficientData},
		EndpointURL:      "http://foobar",
		ExpectedHTTPCode: 200,
		MaxResponseTime:  1000,
		Interval:         60,
		EmailAlerts:      true,
		Active:           true,
	}
	err = suite.db.Create(testAlarm).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare test server and client
	server, client = testServer(&http.Response{StatusCode: 200})
	defer server.Close()
	suite.service.client = client
	start = time.Now().UTC()
	gorm.NowFunc = func() time.Time {
		return start
	}

	// Just one response time metric will be logged
	suite.mockLogResponseTime(start, testAlarm.ID, nil)

	// Trigger multiple parallel alarm checks with the same watermark
	var (
		concurrency = 4
		errChan     = make(chan error)
	)
	for i := 0; i < concurrency; i++ {
		go suite.alarmCheckWrapper(testAlarm.ID, testAlarm.Watermark.Time, errChan)
		time.Sleep(5 * time.Millisecond)
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
		assert.Equal(suite.T(), ErrCheckAlreadyTriggered, err)
	}

	// Fetch the updated alarm
	alarm = new(Alarm)
	assert.False(suite.T(), suite.service.db.Preload("User").Preload("Incidents").
		First(alarm, testAlarm.ID).RecordNotFound())

	// Watermark updated
	assert.Equal(
		suite.T(),
		start.Format("2006-01-02T15:04:05Z"),
		alarm.Watermark.Time.UTC().Format("2006-01-02T15:04:05Z"),
	)

	// Status OK
	assert.Equal(suite.T(), alarmstates.OK, alarm.AlarmStateID.String)

	// 0 incidents
	assert.Equal(suite.T(), 0, len(alarm.Incidents))

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AlarmsTestSuite) alarmCheckWrapper(alarmID uint, watermark time.Time, errChan chan error) {
	errChan <- suite.service.CheckAlarm(alarmID, watermark)
}
